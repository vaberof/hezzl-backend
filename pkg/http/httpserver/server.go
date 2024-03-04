package httpserver

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AppServer struct {
	Server    *http.Server
	ChiRouter *chi.Mux
	config    *ServerConfig
}

func New(config *ServerConfig) *AppServer {
	chiRouter := chi.NewRouter()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", config.Host, config.Port),
		Handler: chiRouter}

	return &AppServer{
		Server:    httpServer,
		ChiRouter: chiRouter,
		config:    config,
	}
}

func (server *AppServer) StartAsync() <-chan error {
	exitChannel := make(chan error)

	go func() {
		err := server.Server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	return exitChannel
}
