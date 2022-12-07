package main

import (
	"os"

	"golang.org/x/exp/slog"
)

// makeLogger creates a logger with default settings.
// It will print to stdout by default.
// If the environment variable CLOUD_RUN_JOB is set, it will use JSON logging.
// For see structured logging proposal: https://go.googlesource.com/proposal/+/master/design/56345-structured-logging.md
func makeLogger(defaultsAttrs ...slog.Attr) *slog.Logger {
	h := slog.NewTextHandler(os.Stdout)
	h.WithAttrs(defaultsAttrs)

	inCloudRun := os.Getenv("CLOUD_RUN_JOB")
	if inCloudRun != "" {
		jh := slog.NewJSONHandler(os.Stdout)
		return slog.New(jh)
	}

	log := slog.New(h)
	return log
}
