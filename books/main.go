package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/aksbuzz/bookstore-microservice/books/events"
	"github.com/aksbuzz/bookstore-microservice/books/repository"
	"github.com/aksbuzz/bookstore-microservice/books/service"
	"github.com/redis/go-redis/v9"

	"github.com/aksbuzz/bookstore-microservice/shared/config"
	"github.com/aksbuzz/bookstore-microservice/shared/db"
	"github.com/aksbuzz/bookstore-microservice/shared/router"
	"github.com/aksbuzz/bookstore-microservice/shared/server"
)

var (
	// dis
	DSN        string = "host=localhost port=5432 user=postgres password=postgres dbname=postgres sslmode=disable"
	ServerPort uint16 = 3002
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	logger := slog.New(slog.NewJSONHandler(os.Stderr, nil))
	slog.SetDefault(logger)

	router := router.New()
	cfg := config.New(DSN, ServerPort)

	db, err := db.New(cfg.DSN)
	if err != nil {
		cancel()
		slog.Error("failed to connect to db", "error", err.Error())
		return
	}
	repo := repository.New(db)
	rc := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	service := service.New(repo)
	service.Register(router)

	server := server.New(ctx, router, cfg)

	go events.HandleUpdateBestSellers(ctx, rc)

	if err := server.Start(ctx); err != nil {
		cancel()
		slog.Error("failed to start server", "error", err.Error())
		return
	}
}
