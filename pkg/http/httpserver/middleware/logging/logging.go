package logging

import (
	"fmt"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/vaberof/hezzl-backend/pkg/logging/logs"
	"log/slog"
	"net/http"
)

type Middleware struct {
	Handler func(http.Handler) http.Handler
	Logger  *slog.Logger
}

func New(logs *logs.Logs) *Middleware {
	return impl(logs, "")
}

func impl(logs *logs.Logs, serverName string) *Middleware {
	loggerName := "http-server"
	if serverName != "" {
		loggerName = fmt.Sprintf("%s.%s", loggerName, serverName)
	}
	logger := logs.WithName(loggerName)

	handler := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(responseWriter http.ResponseWriter, request *http.Request) {
			path := request.URL.Path
			if path == "" {
				path = "/"
			}
			method := request.Method
			logger.Info("Request started", slog.Group("http", "path", path, "method", method))

			ww := middleware.NewWrapResponseWriter(responseWriter, request.ProtoMajor)

			defer func() {
				status := ww.Status()
				if status == 0 {
					s, ok := request.Context().Value(render.StatusCtxKey).(int)
					if ok && s != status {
						status = s
					}
				}

				if status >= 500 {
					logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "error", "status", status))
					return
				}

				logger.Info("Request finished", slog.Group("http", "path", path, "method", method, "result", "success", "status", status))
			}()

			next.ServeHTTP(ww, request)
		})
	}

	return &Middleware{
		Handler: handler,
		Logger:  logger,
	}
}
