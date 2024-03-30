package events

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aksbuzz/bookstore-microservice/orders/model"
	"github.com/nats-io/nats.go"
)

func OrderPlaced(items []*model.CheckoutItem) error {
	books := make([]uint64, 0, len(items))
	for _, item := range items {
		books = append(books, item.BookId)
	}

	payload, err := json.Marshal(model.OrderPlaced{Books: books})
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		return fmt.Errorf("failed to connect to nats: %w", err)
	}
	defer nc.Close()

	slog.Info("Publishing order.placed event")
	msg, err := nc.Request("order.placed", payload, 5*time.Second)
	if err != nil {
		return fmt.Errorf("failed to request order.placed: %w", err)
	}
	if msg.Data == nil {
		return fmt.Errorf("received nil response from order.placed")
	}

	return nil
}
