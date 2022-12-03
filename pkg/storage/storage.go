package storage

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Store struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InsertOrder(ctx context.Context, userID uint64, total uint64) (int64, error) {
	const query = `INSERT INTO orders(user_id, total) VALUES ($1, $2) RETURNING id;`
	res, err := s.db.ExecContext(ctx, query, userID, total)
	if err != nil {
		return 0, fmt.Errorf("insert order: %w", err)
	}

	orderID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("get order id: %w", err)
	}

	return orderID, nil
}
