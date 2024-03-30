package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aksbuzz/bookstore-microservice/orders/repository"
	"github.com/aksbuzz/bookstore-microservice/orders/service"

	"github.com/aksbuzz/bookstore-microservice/shared/config"
	"github.com/aksbuzz/bookstore-microservice/shared/db"
	"github.com/aksbuzz/bookstore-microservice/shared/router"
	"github.com/aksbuzz/bookstore-microservice/shared/server"
)

var (
	// dis
	DSN        string = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	ServerPort uint16 = 3004
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	cfg := config.New(DSN, ServerPort)
	router := router.New()
	db, err := db.New(cfg.DSN)
	if err != nil {
		cancel()
		slog.Error("failed to connect to db", "error", err.Error())
		return
	}
	repo := repository.New(db)

	service := service.New(repo)
	service.Register(router)

	server := server.New(ctx, router, cfg)

	if err := server.Start(ctx); err != nil {
		cancel()
		slog.Error("failed to start server", "error", err.Error())
		return
	}
}
