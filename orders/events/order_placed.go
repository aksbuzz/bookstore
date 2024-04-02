package events

import (
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/aksbuzz/bookstore-microservice/orders/model"
	"github.com/nats-io/nats.go"
)

func OrderPlaced(nc *nats.Conn, items []*model.CheckoutItem) error {
	var EventSubject = "ORDER.placed"
	bookIDs := make([]uint64, len(items))
	for i, item := range items {
		bookIDs[i] = item.BookId
	}

	payload, err := json.Marshal(model.OrderPlaced{Books: bookIDs})
	if err != nil {
		return fmt.Errorf("failed to marshal order placed payload: %w", err)
	}

	slog.Info(fmt.Sprintf("Publishing '%s' event", EventSubject))
	if err := nc.Publish(EventSubject, payload); err != nil {
		return fmt.Errorf("failed to publish '%s' event: %w", EventSubject, err)
	}
	return nil
}
