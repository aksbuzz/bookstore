package events

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"

	"github.com/aksbuzz/bookstore-microservice/orders/model"
	"github.com/nats-io/nats.go"
)

type bestSeller struct {
	Id       string  `json:"id"`
	Quantity float64 `json:"quantity"`
}

type updateBestsellers struct {
	Books []*bestSeller `json:"books"`
}

func UpdateBestsellers(nc *nats.Conn, items []*model.CheckoutItem) error {
	var EventSubject = "BESTSELLERS.update"

	books := make([]*bestSeller, len(items))
	for i, item := range items {
		books[i] = &bestSeller{
			Id:       strconv.FormatUint(item.BookId, 10),
			Quantity: float64(item.Quantity),
		}
	}

	payload, err := json.Marshal(updateBestsellers{Books: books})
	if err != nil {
		return fmt.Errorf("failed to marshal update bestsellers payload: %w", err)
	}

	slog.Info(fmt.Sprintf("Publishing '%s' event", EventSubject))

	if err := nc.Publish(EventSubject, payload); err != nil {
		return fmt.Errorf("failed to publish '%s' event: %w", EventSubject, err)
	}
	return nil
}
