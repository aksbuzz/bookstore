package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/aksbuzz/bookstore-microservice/orders/events"
	"github.com/aksbuzz/bookstore-microservice/orders/model"
	"github.com/aksbuzz/bookstore-microservice/shared/helper"
	"github.com/go-chi/chi/v5"
)

func (s *Service) ListOrdersWithProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.Repository.List(ctx)
	if err != nil {
		slog.Error("failed to list orders", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	orders, err := json.Marshal(list)
	if err != nil {
		slog.Error("failed to marshal orders", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(orders)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) Checkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var checkout struct {
		Items []*model.CheckoutItem `json:"items"`
	}
	if err := json.NewDecoder(r.Body).Decode(&checkout); err != nil {
		slog.Error("failed to decode checkout", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	now := time.Now().UTC()
	order := &model.CheckoutOrder{
		Date:  now,
		Total: calculateTotal(checkout.Items),
	}
	if err := s.Repository.Checkout(ctx, order, checkout.Items); err != nil {
		slog.Error("failed to checkout", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: replace with NATS streaming
	s.updateBestSellers(ctx, checkout.Items)

	if err := events.OrderPlaced(checkout.Items); err != nil {
		slog.Error("failed to publish event", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Service) GetOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	idParam := chi.URLParam(r, "order_id")
	orderId, err := helper.ParseUint64(idParam)
	if err != nil {
		slog.Error("failed to parse id", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	order, err := s.Repository.Find(ctx, orderId)
	if err != nil {
		slog.Error("failed to find order", "error", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if order == nil {
		slog.Error("order not found")
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(&order); err != nil {
		slog.Error("failed to encode order", "error", err.Error())
		return
	}
}

// TODO: replace with NATS streaming
func (s *Service) updateBestSellers(ctx context.Context, items []*model.CheckoutItem) {
	for _, item := range items {
		_, err := s.Redis.ZIncrBy(ctx, "TopBooks:AllTime", float64(item.Quantity), strconv.FormatUint(item.BookId, 10)).Result()
		if err != nil {
			slog.Error("failed to update bestsellers", "error", err.Error())
			continue
		}
	}
}

func calculateTotal(items []*model.CheckoutItem) float64 {
	total := 0.0
	for _, item := range items {
		total += item.Price * float64(item.Quantity)
	}
	return total
}
