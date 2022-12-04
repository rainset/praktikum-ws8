package storage

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStore_InsertOrder(t *testing.T) {
	db := sqlx.MustOpen("sqlite3", "../../nanomart_test.db")
	defer db.Close()
	s := &Store{db}

	_, err := s.InsertOrder(context.Background(), 42, 100)

	require.NoError(t, err)
	orders := getOrders(db)
	require.Len(t, orders, 1)
	assert.EqualValues(t, 42, orders[0].UserID)
	assert.EqualValues(t, 100, orders[0].Total)
}

type Order struct {
	UserID    uint64    `db:"user_id"`
	Total     uint64    `db:"total"`
	CreatedAt time.Time `db:"created_at"`
}

func getOrders(db *sqlx.DB) []Order {
	const query = `SELECT user_id, total, created_at FROM orders;`

	var res []Order
	if err := db.SelectContext(context.Background(), &res, query); err != nil {
		panic(err)
	}

	return res
}
