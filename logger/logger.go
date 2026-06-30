// Package logger initialises a triple-channel logger:
//   - stdout      — human-readable pretty text format (no ANSI colors)
//   - app.log     — plain text format (no pretty indentation), rolling via lumberjack
//   - error.log   — plain text format (no pretty indentation), rolling via lumberjack
//
// Each channel has an independent log level (supports Off to disable).
// Call Init once after config is loaded.
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

// Init sets up the global slog default logger with three output handlers.
func Init() {
	stdoutLevel := config.ParseLevel(config.Cfg.StdoutLogLevel)
	fileLevel := config.ParseLevel(config.Cfg.FileLogLevel)
	errorLevel := config.ParseLevel(config.Cfg.ErrorLogLevel)

	// Ensure logs directory exists (relative to the executable location)
	exePath, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "fatal: cannot determine executable path: %v\n", err)
		os.Exit(1)
	}
	logsDir := filepath.Join(filepath.Dir(exePath), "logs")

	// If logs exists but is a file, remove it
	info, err := os.Stat(logsDir)
	if err == nil && !info.IsDir() {
		if err := os.Remove(logsDir); err != nil {
			fmt.Fprintf(os.Stderr, "fatal: cannot remove existing logs file: %v\n", err)
			os.Exit(1)
		}
	}
	if err := os.MkdirAll(logsDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "fatal: cannot create logs directory: %v\n", err)
		os.Exit(1)
	}

	var handlers []slog.Handler

	// stdout handler — pretty text format with KV indentation
	if stdoutLevel < config.ParseLevel("Off") {
		stdoutHandler := &prettyTextHandler{
			level:  stdoutLevel,
			writer: os.Stdout,
		}
		handlers = append(handlers, stdoutHandler)
	}

	// app.log handler — plain text format, rolling via lumberjack
	if fileLevel < config.ParseLevel("Off") {
		appWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logsDir, "app.log"),
			MaxSize:    10, // MB
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		}
		appHandler := &plainTextHandler{
			level:  fileLevel,
			writer: appWriter,
		}
		handlers = append(handlers, appHandler)
	}

	// error.log handler — plain text format, rolling via lumberjack
	if errorLevel < config.ParseLevel("Off") {
		errorWriter := &lumberjack.Logger{
			Filename:   filepath.Join(logsDir, "error.log"),
			MaxSize:    10, // MB
			MaxBackups: 3,
			MaxAge:     7, // days
			Compress:   true,
		}
		errorHandler := &plainTextHandler{
			level:  errorLevel,
			writer: errorWriter,
		}
		handlers = append(handlers, errorHandler)
	}

	if len(handlers) == 0 {
		// All channels off — use a no-op handler
		slog.SetDefault(slog.New(noopHandler{}))
		return
	}

	multi := &multiHandler{handlers: handlers}
	slog.SetDefault(slog.New(multi))
}

// ---- prettyTextHandler: stdout only, with KV indentation ----

// prettyTextHandler outputs human-readable logs with KV fields on separate indented lines:
//
//	[2024-01-01 12:00:00] INFO | server started | port=8080
//	                                host=0.0.0.0
type prettyTextHandler struct {
	level  slog.Level
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

func (t *prettyTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= t.level
}

func (t *prettyTextHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := r.Level.String()

	var kvs []string
	for _, a := range t.attrs {
		kvs = append(kvs, formatAttr(a))
	}
	r.Attrs(func(a slog.Attr) bool {
		kvs = append(kvs, formatAttr(a))
		return true
	})

	line := fmt.Sprintf("[%s] %s | %s", timeStr, levelStr, r.Message)

	if len(kvs) > 0 {
		line += " | " + kvs[0]
		prefix := strings.Repeat(" ", len(fmt.Sprintf("[%s] %s | ", timeStr, levelStr)))
		for _, kv := range kvs[1:] {
			line += "\n" + prefix + kv
		}
	}

	fmt.Fprintln(t.writer, line)
	return nil
}

func (t *prettyTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(t.attrs)+len(attrs))
	copy(newAttrs, t.attrs)
	copy(newAttrs[len(t.attrs):], attrs)
	return &prettyTextHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  newAttrs,
		groups: t.groups,
	}
}

func (t *prettyTextHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(t.groups)+1)
	copy(newGroups, t.groups)
	newGroups[len(t.groups)] = name
	return &prettyTextHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  t.attrs,
		groups: newGroups,
	}
}

// ---- plainTextHandler: file output, single-line plain text ----

// plainTextHandler outputs single-line plain text logs (no pretty indentation):
//
//	[2024-01-01 12:00:00] INFO | server started | port=8080 host=0.0.0.0
type plainTextHandler struct {
	level  slog.Level
	writer io.Writer
	attrs  []slog.Attr
	groups []string
}

func (t *plainTextHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= t.level
}

func (t *plainTextHandler) Handle(_ context.Context, r slog.Record) error {
	timeStr := r.Time.Format("2006-01-02 15:04:05")
	levelStr := r.Level.String()

	var kvs []string
	for _, a := range t.attrs {
		kvs = append(kvs, formatAttr(a))
	}
	r.Attrs(func(a slog.Attr) bool {
		kvs = append(kvs, formatAttr(a))
		return true
	})

	line := fmt.Sprintf("[%s] %s | %s", timeStr, levelStr, r.Message)

	if len(kvs) > 0 {
		line += " | " + strings.Join(kvs, " ")
	}

	fmt.Fprintln(t.writer, line)
	return nil
}

func (t *plainTextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	newAttrs := make([]slog.Attr, len(t.attrs)+len(attrs))
	copy(newAttrs, t.attrs)
	copy(newAttrs[len(t.attrs):], attrs)
	return &plainTextHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  newAttrs,
		groups: t.groups,
	}
}

func (t *plainTextHandler) WithGroup(name string) slog.Handler {
	newGroups := make([]string, len(t.groups)+1)
	copy(newGroups, t.groups)
	newGroups[len(t.groups)] = name
	return &plainTextHandler{
		level:  t.level,
		writer: t.writer,
		attrs:  t.attrs,
		groups: newGroups,
	}
}

// ---- shared helpers ----

// formatAttr formats a single attribute as "key=value".
func formatAttr(a slog.Attr) string {
	return a.Key + "=" + a.Value.String()
}

// ---- multiHandler: fan-out ----

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

// ---- noopHandler: all channels off ----

type noopHandler struct{}

func (noopHandler) Enabled(context.Context, slog.Level) bool  { return false }
func (noopHandler) Handle(context.Context, slog.Record) error { return nil }
func (n noopHandler) WithAttrs([]slog.Attr) slog.Handler      { return n }
func (n noopHandler) WithGroup(string) slog.Handler           { return n }

// Sync flushes any buffered log output. Call before program exit.
func Sync() {
	_, _ = io.WriteString(io.Discard, "")
}
