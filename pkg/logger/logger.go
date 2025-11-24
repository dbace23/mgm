package logger

import (
	"log/slog"
	"os"
)

var log *slog.Logger

func Init(env string) {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	if env == "development" {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, opts)
	log = slog.New(handler)
}

func Info(msg string, args ...any) {
	log.Info(msg, args...)
}

func Error(msg string, args ...any) {
	log.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	log.Debug(msg, args...)
}

func Warn(msg string, args ...any) {
	log.Warn(msg, args...)
}

func Fatal(msg string, args ...any) {
	log.Error(msg, args...)
	os.Exit(1)
}
