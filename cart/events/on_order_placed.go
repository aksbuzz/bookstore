package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/aksbuzz/bookstore-microservice/cart/model"
	"github.com/aksbuzz/bookstore-microservice/cart/repository"
	"github.com/nats-io/nats.go"
)

func HandleOrderPlaced(ctx context.Context, repo repository.CartRepository) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		slog.Error("failed to connect to nats", "error", err)
		return
	}
	defer nc.Close()

	slog.Info("Subscribing to order.placed event")
	_, err = nc.Subscribe("order.placed", func(msg *nats.Msg) {
		payload := &model.OrderPlaced{}
		if err := json.Unmarshal(msg.Data, payload); err != nil {
			msg.Respond(nil)
			return
		}

		for _, bookId := range payload.Books {
			if err := repo.RemoveFromCart(ctx, bookId); err != nil {
				msg.Respond(nil)
				return
			}
		}

		msg.Respond([]byte("OK"))
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
