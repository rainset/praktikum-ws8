package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type storageStub struct{}

func (storageStub) InsertOrder(ctx context.Context, userID uint64, total uint64) (int64, error) {
	return 1, nil
}

type storageMock struct {
	insertFunc func(ctx context.Context, userID uint64, total uint64) (int64, error)
}

func (m *storageMock) InsertOrder(ctx context.Context, userID uint64, total uint64) (int64, error) {
	return m.insertFunc(ctx, userID, total)
}

func TestApp_HandleCreateOrder(t *testing.T) {
	storageMock := &storageMock{}

	app := &App{store: storageMock}

	req := `{
	"user_id": 42,
	"items": [
		{
			"name": "Milk",
			"price": 45
		},
		{
			"name": "Bread",
			"price": 12
		}
	]
}`

	w := httptest.NewRecorder()
	r := httptest.NewRequest("POST", "localhost:5301/api/v1/order", bytes.NewBuffer([]byte(req)))

	storageMock.insertFunc = func(ctx context.Context, userID, total uint64) (int64, error) {
		return 3, nil
	}

	app.HandleCreateOrder(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	res := make(map[string]interface{})
	json.NewDecoder(w.Result().Body).Decode(&res)
	assert.EqualValues(t, 3, res["order_id"])
}
