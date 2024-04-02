package events

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

type bestSeller struct {
	Id       string  `json:"id"`
	Quantity float64 `json:"quantity"`
}

type updateBestsellers struct {
	Books []*bestSeller `json:"books"`
}

func HandleUpdateBestSellers(ctx context.Context, rc *redis.Client) {
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		slog.Error("failed to connect to nats", "error", err)
		return
	}
	defer nc.Close()

	slog.Info("Subscribing to BESTSELLERS.update event")
	_, err = nc.Subscribe("BESTSELLERS.update", func(msg *nats.Msg) {
		payload := &updateBestsellers{}
		if err := json.Unmarshal(msg.Data, payload); err != nil {
			slog.Error("failed to unmarshal update bestsellers payload", "error", err)
			return
		}

		for _, book := range payload.Books {
			_, err := rc.ZIncrBy(ctx, "TopBooks:AllTime", book.Quantity, book.Id).Result()
			if err != nil {
				slog.Error("failed to update bestsellers", "error", err.Error())
				continue
			}
		}

		slog.Info("Updated bestsellers")
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
