package service

import (
	"context"
	"log/slog"
	"net"
	"net/http"
)

type Service struct {
	// Shared memory for the service
	httpserver *http.Server
	addr       string
}

func New(ctx context.Context, addr string) *Service {
	return &Service{
		httpserver: &http.Server{
			Addr:        addr,
			BaseContext: func(net.Listener) context.Context { return ctx },
		},
		addr: addr,
	}
}

func (s *Service) Run() error {
	http.HandleFunc("/", index)

	slog.Info("starting service", "addr", s.addr)
	if err := s.httpserver.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (s *Service) Close(ctx context.Context) error {
	return s.httpserver.Shutdown(ctx)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
