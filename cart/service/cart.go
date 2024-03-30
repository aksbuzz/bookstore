package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/aksbuzz/bookstore-microservice/cart/model"
	"github.com/aksbuzz/bookstore-microservice/shared/helper"
	"github.com/go-chi/chi/v5"
)

func (s *Service) ListCartItems(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.Repository.List(ctx)
	if err != nil {
		slog.Error("failed to list cart items", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	items, err := json.Marshal(list)
	if err != nil {
		slog.Error("failed to marshal cart items", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(items)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	item := model.AddCart{}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		slog.Error("failed to decode item", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	exist, err := s.Repository.Find(ctx, item.BookId)
	if err != nil {
		slog.Error("failed to find item", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if exist != nil {
		item.Quantity = exist.Quantity + 1
		if err := s.Repository.UpdateCart(ctx, item.BookId, item.Quantity); err != nil {
			slog.Error("failed to update item", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := s.Repository.AddToCart(ctx, &item); err != nil {
			slog.Error("failed to add item", "error", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Service) UpdateItemQuantity(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "book_id")
	bookId, err := helper.ParseUint64(idParam)
	if err != nil {
		slog.Error("failed to parse id", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	update := model.UpdateCart{}
	if err := json.NewDecoder(r.Body).Decode(&update); err != nil {
		slog.Error("failed to decode update", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Repository.UpdateCart(ctx, bookId, update.Quantity); err != nil {
		slog.Error("failed to update item", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "book_id")
	bookId, err := helper.ParseUint64(idParam)
	if err != nil {
		slog.Error("failed to parse id", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := s.Repository.RemoveFromCart(ctx, bookId); err != nil {
		slog.Error("failed to remove item", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetCartItem(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "book_id")
	bookId, err := helper.ParseUint64(idParam)
	if err != nil {
		slog.Error("failed to parse id", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	item, err := s.Repository.Find(ctx, bookId)
	if err != nil {
		slog.Error("failed to find item", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if item == nil {
		slog.Error("item not found in cart")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&item); err != nil {
		slog.Error("failed to encode item", "error", err.Error())
		return
	}
}
