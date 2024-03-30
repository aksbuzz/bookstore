package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aksbuzz/bookstore-microservice/books/model"
)

type BookRepository interface {
	List(ctx context.Context) ([]*model.Book, error)
	Find(ctx context.Context, id uint64) (*model.Book, error)
}

func (r *Repository) List(ctx context.Context) ([]*model.Book, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, name, author, category, price, rating, cover
		FROM book
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		var book model.Book
		if err := rows.Scan(
			&book.Id,
			&book.Name,
			&book.Author,
			&book.Category,
			&book.Price,
			&book.Rating,
			&book.Cover,
		); err != nil {
			return nil, err
		}
		books = append(books, &book)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return books, nil
}

func (r *Repository) Find(ctx context.Context, id uint64) (*model.Book, error) {
	var book model.Book

	err := r.db.QueryRowContext(ctx, `
		SELECT id, name, author, category, price, rating, cover
		FROM book
		WHERE id = $1
		`, id).
		Scan(&book.Id,
			&book.Name,
			&book.Author,
			&book.Category,
			&book.Price,
			&book.Rating,
			&book.Cover)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &book, nil
}
