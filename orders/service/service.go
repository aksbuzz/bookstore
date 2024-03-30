package service

import (
	"github.com/aksbuzz/bookstore-microservice/orders/repository"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repository repository.OrderRepository
	Redis      *redis.Client
}

func New(repository repository.OrderRepository) *Service {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return &Service{
		Repository: repository,
		Redis:      redisClient,
	}
}

func (s *Service) Register(router *chi.Mux) {
	router.Route("/orders", func(r chi.Router) {
		r.Get("/", s.ListOrdersWithProducts)
		r.Post("/", s.Checkout)
		r.Get("/{order_id}", s.GetOrder)
	})
}
