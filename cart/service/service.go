package service

import (
	"github.com/aksbuzz/bookstore-microservice/cart/repository"
	"github.com/go-chi/chi/v5"
)

type Service struct {
	Repository repository.CartRepository
}

func New(repository repository.CartRepository) *Service {
	return &Service{
		Repository: repository,
	}
}

func (s *Service) Register(router *chi.Mux) {
	router.Route("/cart", func(r chi.Router) {
		r.Get("/", s.ListCartItems)
		r.Post("/", s.AddItemToCart)
		r.Put("/{book_id}", s.UpdateItemQuantity)
		r.Delete("/{book_id}", s.RemoveItemFromCart)
		r.Get("/{book_id}", s.GetCartItem)
	})
}
