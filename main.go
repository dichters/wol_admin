package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"wol_admin/antishake"
	"wol_admin/config"
	"wol_admin/handler"
	"wol_admin/logger"
	"wol_admin/version"
)

//go:embed dist/*
var staticFS embed.FS

func main() {
	// Handle `./wol_admin version` subcommand
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("wol_admin %s %s %s\n", version.Version, version.Arch, version.BuildTime)
		return
	}

	// 1. Load config
	config.Load(".")

	// 2. Init logger (dual-channel, independent levels)
	logger.Init()
	slog.Info("config loaded", "port", config.Cfg.ServerPort)

	// 3. Init anti-shake locker (Redis when enabled, skip otherwise)
	locker := antishake.New()
	defer locker.Close()

	// 4. Build HTTP mux
	mux := http.NewServeMux()

	// API routes under /wol/api/
	apiHandler := handler.NewAPIHandler(locker)
	mux.HandleFunc("POST /wol/api/wol", apiHandler.WOL)
	mux.HandleFunc("POST /wol/api/shutdown", apiHandler.Shutdown)
	mux.HandleFunc("GET /wol/api/version", apiHandler.Version)

	// Static files — serve embedded Vue3 frontend under /wol/
	distFS, err := fs.Sub(staticFS, "dist")
	if err != nil {
		slog.Error("failed to read embedded static files", "error", err)
		os.Exit(1)
	}
	fileServer := http.FileServer(http.FS(distFS))
	spaHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Strip the /wol prefix so fileServer looks in dist/
		r.URL.Path = strings.TrimPrefix(r.URL.Path, "/wol")
		if r.URL.Path == "" {
			r.URL.Path = "/"
		}

		// Try to serve the static file first
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Check if file exists in embedded FS
		if _, err := fs.Stat(distFS, strings.TrimPrefix(path, "/")); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// SPA fallback: serve index.html for client-side routes (e.g. /en)
		r.URL.Path = "/index.html"
		fileServer.ServeHTTP(w, r)
	})
	mux.Handle("GET /wol/", spaHandler)

	// Redirect root to /wol/
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/wol/", http.StatusFound)
			return
		}
		http.NotFound(w, r)
	})

	// 5. Start server
	addr := "0.0.0.0:" + config.Cfg.ServerPort
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	go func() {
		slog.Info("server starting", "addr", addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server listen failed", "error", err)
			os.Exit(1)
		}
	}()

	// 6. Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("shutting down server...")
	server.Close()
	logger.Sync()
	slog.Info("server stopped")
}
