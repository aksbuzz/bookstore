package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aksbuzz/bookstore-microservice/cart/events"
	"github.com/aksbuzz/bookstore-microservice/cart/repository"
	"github.com/aksbuzz/bookstore-microservice/cart/service"
	"github.com/nats-io/nats.go"

	"github.com/aksbuzz/bookstore-microservice/shared/config"
	"github.com/aksbuzz/bookstore-microservice/shared/db"
	"github.com/aksbuzz/bookstore-microservice/shared/router"
	"github.com/aksbuzz/bookstore-microservice/shared/server"
)

var (
	// dis
	DSN        string = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	ServerPort uint16 = 3003
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize logger
	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	// Initialize router
	router := router.New()

	// Initialize config
	cfg := config.New(DSN, ServerPort)

	// Initialize database
	db, err := db.New(cfg.DSN)
	if err != nil {
		slog.Error("failed to connect to db", "error", err.Error())
		return
	}
	defer db.Close()

	// Initialize repository
	repo := repository.New(db)

	// Initialize NATS
	nc, err := nats.Connect("nats://localhost:4222")
	if err != nil {
		slog.Error("failed to connect to nats", "error", err.Error())
		return
	}
	defer nc.Close()

	// Initialize service
	service := service.New(repo)
	service.Register(router)

	go events.HandleOrderPlaced(ctx, nc, repo)

	// Initialize server
	server := server.New(ctx, router, cfg)

	if err := server.Start(ctx); err != nil {
		slog.Error("failed to start server", "error", err.Error())
		return
	}
}
