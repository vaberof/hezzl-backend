package logs

import (
	"io"
	"log/slog"
)

type Logs struct {
	root *slog.Logger
}

func New(w io.Writer, opts *slog.HandlerOptions) *Logs {
	jsonLogger := getJSONLogger(w, opts)
	return &Logs{root: jsonLogger}
}

func (logs *Logs) WithName(name string) *slog.Logger {
	return logs.root.With("logger-name", name)
}

func (logs *Logs) GetLogger() *slog.Logger {
	return logs.root
}

func getJSONLogger(w io.Writer, opts *slog.HandlerOptions) *slog.Logger {
	jsonHandler := slog.NewJSONHandler(w, opts)
	logger := slog.New(jsonHandler)
	return logger
}
