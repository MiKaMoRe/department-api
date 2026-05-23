package logger

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"department-api/internal/config"
)

type teeHandler struct {
	a, b slog.Handler
}

func (t *teeHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return t.a.Enabled(ctx, level)
}

func (t *teeHandler) Handle(ctx context.Context, r slog.Record) error {
	err1 := t.a.Handle(ctx, r)
	err2 := t.b.Handle(ctx, r.Clone())
	return errors.Join(err1, err2)
}

func (t *teeHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &teeHandler{a: t.a.WithAttrs(attrs), b: t.b.WithAttrs(attrs)}
}

func (t *teeHandler) WithGroup(name string) slog.Handler {
	return &teeHandler{a: t.a.WithGroup(name), b: t.b.WithGroup(name)}
}

const defaultDevLogFile = "logs/dev.log"

var Log *slog.Logger

// devLogFile is set in dev mode when logging to a file alongside stdout; Close releases it.
var devLogFile *os.File

func MustInit(env string) {
	var handler slog.Handler

	switch env {
	case config.EnvDev:
		stdoutLog := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug})
		openDevLogFile()
		if devLogFile != nil {
			fileLog := slog.NewJSONHandler(devLogFile, &slog.HandlerOptions{Level: slog.LevelDebug})
			handler = &teeHandler{a: stdoutLog, b: fileLog}
		} else {
			handler = stdoutLog
		}
	case config.EnvTest:
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	case config.EnvProd:
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})
	default:
		panic("Logger: invalid environment")
	}

	Log = slog.New(handler)

	slog.SetDefault(Log)

	Log.Info("Logger initialized", "ENV", env)
}

// Close releases the dev log file if one was opened. Safe to call multiple times.
func Close() {
	if devLogFile != nil {
		_ = devLogFile.Close()
		devLogFile = nil
	}
}

func openDevLogFile() {
	path := os.Getenv("LOG_FILE_PATH")
	if path == "" {
		path = defaultDevLogFile
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "logger: cannot create log directory %q: %v (using stdout only)\n", dir, err)
		return
	}

	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "logger: cannot open log file %q: %v (using stdout only)\n", path, err)
		return
	}

	devLogFile = f
}

func Debug(msg string, args ...any) {
	Log.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	Log.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	Log.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	Log.Error(msg, args...)
}
