package service

import (
	"github.com/aksbuzz/bookstore-microservice/books/repository"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
)

type Service struct {
	Repository repository.BookRepository
	Redis      *redis.Client
}

func New(repository repository.BookRepository) *Service {
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
	router.Route("/books", func(r chi.Router) {
		r.Get("/", s.ListBooks)
		r.Get("/bestsellers", s.ListBestSellers)
		r.Get("/{id}", s.GetBook)
	})
}
