// Package logger initialises a dual-channel text logger:
//   - stdout  — human-readable text format (no ANSI colors)
//   - file    — human-readable text format, rolling via lumberjack
//
// Call Init once after config is loaded. The global default logger is replaced
// so all slog.Info / slog.Error etc. calls in the program use both channels.
package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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

	// stdout handler — human-readable text format, independent level
	stdoutHandler := &textHandler{
		level:  stdoutLevel,
		writer: os.Stdout,
	}

	// file handler — text format, rolling via lumberjack, independent level
	fileWriter := &lumberjack.Logger{
		Filename:   filepath.Join("logs", "app.log"),
		MaxSize:    10, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
		Compress:   true,
	}

	fileHandler := &textHandler{
		level:  fileLevel,
		writer: fileWriter,
	}

	// Multi-handler: fan-out to both channels
	multi := &multiHandler{
		handlers: []slog.Handler{stdoutHandler, fileHandler},
	}

	slog.SetDefault(slog.New(multi))
}

// textHandler outputs human-readable logs to stdout:
//
//	[2024-01-01 12:00:00] INFO | server started | port=8080
//	                                host=0.0.0.0
type textHandler struct {
	level  slog.Level
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

func (t *textHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= t.level
}

func (t *textHandler) Handle(_ context.Context, r slog.Record) error {
	// Time format: 2006-01-02 15:04:05
	timeStr := r.Time.Format("2006-01-02 15:04:05")

	// Level string
	levelStr := r.Level.String()

	// Collect all KV pairs (attrs from context + record attrs)
	var kvs []string
	for _, a := range t.attrs {
		kvs = append(kvs, formatAttr(a))
	}
	r.Attrs(func(a slog.Attr) bool {
		kvs = append(kvs, formatAttr(a))
		return true
	})

	// Build output line
	line := fmt.Sprintf("[%s] %s | %s", timeStr, levelStr, r.Message)

	if len(kvs) > 0 {
		line += " | " + kvs[0]
		// Subsequent KV fields on new lines with aligned indentation
		prefix := strings.Repeat(" ", len(fmt.Sprintf("[%s] %s | ", timeStr, levelStr)))
		for _, kv := range kvs[1:] {
			line += "\n" + prefix + kv
		}
	}

	fmt.Fprintln(t.writer, line)
	return nil
}

func (t *textHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(t.attrs)+len(attrs))
	copy(newAttrs, t.attrs)
	copy(newAttrs[len(t.attrs):], attrs)
	return &textHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  newAttrs,
		groups: t.groups,
	}
}

func (t *textHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(t.groups)+1)
	copy(newGroups, t.groups)
	newGroups[len(t.groups)] = name
	return &textHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  t.attrs,
		groups: newGroups,
	}
}

// formatAttr formats a single attribute as "key=value".
func formatAttr(a slog.Attr) string {
	if a.Value.Kind() == slog.KindString {
		return a.Key + "=" + a.Value.String()
	}
	return a.Key + "=" + a.Value.String()
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
	_, _ = io.WriteString(io.Discard, "")
}
