package repository

import (
	"context"

	"github.com/aksbuzz/bookstore-microservice/orders/model"
)

type OrderRepository interface {
	List(ctx context.Context) ([]*model.Order, error)
	Checkout(ctx context.Context, order *model.CheckoutOrder, items []*model.CheckoutItem) error
	Find(ctx context.Context, id uint64) (*model.Order, error)
}

func (r *Repository) List(ctx context.Context) ([]*model.Order, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT o.id, o.date, o.total, b.author, b.name, b.cover, oi.id, oi.quantity, oi.price
		FROM "order" o
		JOIN order_item oi ON o.id = oi.order_id
		JOIN book b ON oi.book_id = b.id
		ORDER BY o.date
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := make(map[uint64]*model.Order)
	for rows.Next() {
		var order model.Order
		var item model.OrderItem
		if err := rows.Scan(
			&order.Id,
			&order.Date,
			&order.Total,
			&item.Author,
			&item.Name,
			&item.Cover,
			&item.Id,
			&item.Quantity,
			&item.Price,
		); err != nil {
			return nil, err
		}

		if _, ok := orders[order.Id]; !ok {
			orders[order.Id] = &order
		}
		orders[order.Id].Items = append(orders[order.Id].Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	var result []*model.Order
	for _, order := range orders {
		result = append(result, order)
	}

	return result, nil
}

func (r *Repository) Checkout(ctx context.Context, order *model.CheckoutOrder, items []*model.CheckoutItem) error {
	txn, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer txn.Rollback()

	var orderId uint64
	if err := txn.QueryRowContext(ctx, `
		INSERT INTO "order" (date, total)
		VALUES ($1, $2)
		RETURNING id
	`, order.Date, order.Total).Scan(&orderId); err != nil {
		return err
	}

	for _, item := range items {
		if _, err := txn.ExecContext(ctx, `
			INSERT INTO order_item (order_id, book_id, quantity, price)
			VALUES ($1, $2, $3, $4)
		`, orderId, item.BookId, item.Quantity, item.Price); err != nil {
			return err
		}
	}

	if err := txn.Commit(); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Find(ctx context.Context, id uint64) (*model.Order, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT o.id, o.date, o.total
		FROM "order" o
		WHERE o.id = $1
	`, id)

	var order model.Order
	if err := row.Scan(&order.Id, &order.Date, &order.Total); err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT b.author, b.name, b.cover, oi.quantity, oi.price
		FROM order_item oi
		JOIN book b ON oi.book_id = b.id
		WHERE oi.order_id = $1
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var item model.OrderItem
		if err := rows.Scan(
			&item.Author,
			&item.Name,
			&item.Cover,
			&item.Quantity,
			&item.Price,
		); err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &order, nil
}
