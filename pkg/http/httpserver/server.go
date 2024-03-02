package httpserver

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type AppServer struct {
	Server  *chi.Mux
	config  *ServerConfig
	address string
}

func New(config *ServerConfig) *AppServer {
	chiServer := chi.NewRouter()

	return &AppServer{
		Server:  chiServer,
		config:  config,
		address: fmt.Sprintf("%s:%d", config.Host, config.Port),
	}
}

func (server *AppServer) StartAsync() <-chan error {
	exitChannel := make(chan error)

	go func() {
		err := http.ListenAndServe(server.address, server.Server)
		if !errors.Is(err, http.ErrServerClosed) {
			exitChannel <- err
			return
		} else {
			exitChannel <- nil
		}
	}()

	return exitChannel
}
