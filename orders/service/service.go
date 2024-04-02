package service

import (
	"github.com/aksbuzz/bookstore-microservice/orders/repository"
	"github.com/go-chi/chi/v5"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repository repository.OrderRepository
	Redis      *redis.Client
	nats       *nats.Conn
}

func New(repository repository.OrderRepository, rc *redis.Client, nc *nats.Conn) *Service {
	return &Service{
		Repository: repository,
		Redis:      rc,
		nats:       nc,
	}
}

func (s *Service) Register(router *chi.Mux) {
	router.Route("/orders", func(r chi.Router) {
		r.Get("/", s.ListOrdersWithProducts)
		r.Post("/", s.Checkout)
		r.Get("/{order_id}", s.GetOrder)
	})
}
