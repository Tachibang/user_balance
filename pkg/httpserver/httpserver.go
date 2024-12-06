package httpserver

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultShutdownTimeout = 15 * time.Second
	defaultAddr            = ":8080"
)

type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
}

func New(handler http.Handler, opts ...Option) *Server {
	s := &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
			Addr:         defaultAddr,
		},
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
	}

	for _, opt := range opts {
		opt(s)
	}

	s.start()
	return s
}

func (s *Server) start() {
	go func() {
		err := s.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("Ошибка при запуске сервера: %v\n", err)
		}
		s.notify <- err
		close(s.notify)
	}()
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		fmt.Printf("Ошибка при остановке сервера: %v\n", err)
	}
	return err
}
