package api

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type storageStub struct{}

func (storageStub) InsertOrder(ctx context.Context, userID uint64, total uint64) (int64, error) {
	return 1, nil
}

func TestApp_HandleCreateOrder(t *testing.T) {
	app := &App{store: &storageStub{}}

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

	app.HandleCreateOrder(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)
}
