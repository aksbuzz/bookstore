package repository

import (
	"database/sql"

	"github.com/aksbuzz/bookstore-microservice/shared/db"
)

type Repository struct {
	db *sql.DB
}

func New(db *db.DB) BookRepository {
	return &Repository{
		db: db.Instance(),
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
