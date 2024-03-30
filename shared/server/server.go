package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aksbuzz/bookstore-microservice/shared/config"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	httpServer *http.Server
	Config     *config.Config
}

func New(ctx context.Context, router *chi.Mux, config *config.Config) *Server {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.ServerPort),
		Handler: router,
	}

	s := &Server{
		httpServer: httpServer,
		Config:     config,
	}

	return s
}

func (s *Server) Start(ctx context.Context) error {
	fmt.Println("Starting server on port", s.Config.ServerPort)
	if err := s.httpServer.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}

	return nil
}

func (s *Server) Shutdown(ctx context.Context) error {
	timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := s.httpServer.Shutdown(timeout); err != nil {
		return fmt.Errorf("failed to shutdown server: %w", err)
	}

	return nil
}
