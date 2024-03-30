package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/aksbuzz/bookstore-microservice/cart/model"
)

type CartRepository interface {
	List(ctx context.Context) ([]*model.CartItem, error)
	AddToCart(ctx context.Context, item *model.AddCart) error
	UpdateCart(ctx context.Context, bookId uint64, quantity uint8) error
	RemoveFromCart(ctx context.Context, bookId uint64) error
	Find(ctx context.Context, bookId uint64) (*model.CartItem, error)
}

func (r *Repository) List(ctx context.Context) ([]*model.CartItem, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT c.id, c.book_id, c.quantity, c.price, b.author, b.name, b.cover
		FROM cart c
		JOIN book b ON c.book_id = b.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.CartItem
	for rows.Next() {
		var item model.CartItem
		if err := rows.Scan(
			&item.Id,
			&item.BookId,
			&item.Quantity,
			&item.Price,
			&item.Author,
			&item.Name,
			&item.Cover,
		); err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func (r *Repository) AddToCart(ctx context.Context, item *model.AddCart) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO cart (book_id, quantity, price)
		VALUES ($1, $2, $3)
	`, item.BookId, item.Quantity, item.Price)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateCart(ctx context.Context, bookId uint64, quantity uint8) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE cart
		SET quantity = $1
		WHERE book_id = $2
	`, quantity, bookId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) RemoveFromCart(ctx context.Context, bookId uint64) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM cart
		WHERE book_id = $1
	`, bookId)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Find(ctx context.Context, bookId uint64) (*model.CartItem, error) {
	var item model.CartItem

	err := r.db.QueryRowContext(ctx, `
		SELECT c.id, c.book_id, c.quantity, c.price, b.author, b.name, b.cover
		FROM cart c
		JOIN book b ON c.book_id = b.id
		WHERE book_id = $1
	`, bookId).Scan(
		&item.Id,
		&item.BookId,
		&item.Quantity,
		&item.Price,
		&item.Author,
		&item.Name,
		&item.Cover)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &item, nil
}
