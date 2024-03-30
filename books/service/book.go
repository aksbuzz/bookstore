package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/aksbuzz/bookstore-microservice/books/model"
	"github.com/aksbuzz/bookstore-microservice/shared/helper"
	"github.com/go-chi/chi/v5"
)

func (s *Service) GetBook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	idParam := chi.URLParam(r, "id")
	bookId, err := helper.ParseUint64(idParam)
	if err != nil {
		slog.Error("failed to parse id", "error", err.Error())
		return
	}

	book, err := s.Repository.Find(ctx, bookId)
	if err != nil {
		slog.Error("failed to find book", "error", err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(&book); err != nil {
		slog.Error("failed to encode book", "error", err.Error())
		return
	}
}

func (s *Service) ListBooks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	list, err := s.Repository.List(ctx)
	if err != nil {
		slog.Error("failed to list books", "error", err.Error())
		return
	}

	books, err := json.Marshal(list)
	if err != nil {
		slog.Error("failed to marshal books", "error", err.Error())
		return
	}

	w.Write(books)
	w.WriteHeader(http.StatusOK)
}

func (s *Service) ListBestSellers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	members, err := s.Redis.ZRevRangeWithScores(ctx, "TopBooks:AllTime", 0, 9).Result()
	if err != nil {
		slog.Error("failed to get top books", "error", err.Error())
		return
	}

	// slice of bestsellers
	bestsellers := []model.BestSeller{}
	for _, member := range members {
		book := model.BestSeller{Sales: member.Score}

		bookId, err := strconv.ParseUint(member.Member.(string), 10, 64)
		if err != nil {
			slog.Error("failed to parse book id", "error", err.Error())
			continue
		}
		book.Id = bookId
		bestsellers = append(bestsellers, book)
	}

	if err := json.NewEncoder(w).Encode(&bestsellers); err != nil {
		slog.Error("failed to encode bestsellers", "error", err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}
