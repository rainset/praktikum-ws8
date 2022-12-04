package storage

import (
	"context"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type StorageTestSuite struct {
	suite.Suite
	db *sqlx.DB
}

func (suite *StorageTestSuite) SetupTest() {
	suite.db = sqlx.MustOpen("sqlite3", "../../nanomart_test.db")
	goose.SetDialect("sqlite3")
	require.NoError(suite.T(), goose.Up(suite.db.DB, "../../migrations"))
}

func (suite *StorageTestSuite) TearDownTest() {
	suite.db.ExecContext(context.Background(), `DELETE FROM orders;`)
}

func (suite *StorageTestSuite) TestInsertOrder() {
	s := &Store{suite.db}

	_, err := s.InsertOrder(context.Background(), 42, 100)

	suite.Require().NoError(err)
	orders := getOrders(suite.db)
	suite.Require().Len(orders, 1)
	suite.EqualValues(42, orders[0].UserID)
	suite.EqualValues(100, orders[0].Total)
}

func TestStorageTestSuite(t *testing.T) {
	suite.Run(t, new(StorageTestSuite))
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
