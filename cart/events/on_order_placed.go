package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aksbuzz/bookstore-microservice/cart/model"
	"github.com/aksbuzz/bookstore-microservice/cart/repository"
	"github.com/nats-io/nats.go"
)

func HandleOrderPlaced(ctx context.Context, nc *nats.Conn, repo repository.CartRepository) {
	slog.Info("Subscribing to order.placed event")
	_, err := nc.Subscribe("ORDER.placed", func(msg *nats.Msg) {
		payload := &model.OrderPlaced{}
		if err := json.Unmarshal(msg.Data, payload); err != nil {
			slog.Error("failed to unmarshal order placed payload", "error", err)
			return
		}

		for _, bookId := range payload.Books {
			if err := repo.RemoveFromCart(ctx, bookId); err != nil {
				slog.Error("failed to remove book from cart", "error", err)
				return
			}
		}

		slog.Info("Deleted items from cart")
	})
	if err != nil {
		slog.Error("failed to subscribe to order.placed event", "error", err)
		return
	}

	nc.Flush()
	if err := nc.LastError(); err != nil {
		slog.Error("failed to flush subscription", "error", err)
		return
	}

	// keep the subscription alive
	select {}
}
