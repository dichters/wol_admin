// Package logger initialises a dual-channel structured logger:
//   - stdout  (console, independent level)
//   - file    (lumberjack rolling, independent level)
//
// Call Init once after config is loaded. The global default logger is replaced
// so all slog.Info / slog.Error etc. calls in the program use both channels.
package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"wol_admin/config"

	"gopkg.in/natefinch/lumberjack.v2"
)

// Init sets up the global slog default logger with dual output handlers.
func Init() {
	stdoutLevel := config.ParseLevel(config.Cfg.StdoutLogLevel)
	fileLevel := config.ParseLevel(config.Cfg.FileLogLevel)

	// Ensure logs directory exists
	if err := os.MkdirAll("logs", 0o755); err != nil {
		slog.Error("failed to create logs directory", "error", err)
		os.Exit(1)
	}

	// stdout handler — JSON, independent level
	stdoutHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: stdoutLevel,
	})

	// file handler — JSON, rolling via lumberjack, independent level
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}

	fileHandler := slog.NewJSONHandler(fileWriter, &slog.HandlerOptions{
		Level: fileLevel,
	})

	// Multi-handler: fan-out to both channels
	multi := &multiHandler{
		handlers: []slog.Handler{stdoutHandler, fileHandler},
	}

	slog.SetDefault(slog.New(multi))
}

// multiHandler fans out every log record to multiple handlers.
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: newHandlers}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	newHandlers := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		newHandlers[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: newHandlers}
}

// Sync flushes any buffered log output. Call before program exit.
func Sync() {
	// lumberjack doesn't buffer, but be safe: close the file writer
	// via a no-op — it flushes on each Write.
	_, _ = io.WriteString(io.Discard, "")
}
