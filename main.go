package main

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"wol_admin/antishake"
	"wol_admin/config"
	"wol_admin/handler"
	"wol_admin/logger"
)

//go:embed dist/*
var staticFS embed.FS

func main() {
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

	// API routes — localhost only
	apiHandler := handler.NewAPIHandler(locker)
	mux.HandleFunc("/api/wol", localOnly(apiHandler.WOL))
	mux.HandleFunc("/api/shutdown", localOnly(apiHandler.Shutdown))

	// Static files — serve embedded Vue3 frontend
	distFS, err := fs.Sub(staticFS, "dist")
	if err != nil {
		slog.Error("failed to read embedded static files", "error", err)
		os.Exit(1)
	}
	fileServer := http.FileServer(http.FS(distFS))
	mux.Handle("/", fileServer)

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

// localOnly wraps an http.HandlerFunc to reject non-localhost requests.
func localOnly(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !handler.IsLocalRequest(r) {
			slog.Warn("API access denied: non-local request", "remote", r.RemoteAddr, "path", r.URL.Path)
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
