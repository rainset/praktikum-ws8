package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	mock_api "github.com/kilchik/nanomart/pkg/api/mocks"
)

func TestApp_HandleCreateOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	storageMock := mock_api.NewMockStorage(ctrl)

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

	storageMock.EXPECT().InsertOrder(gomock.Any(), uint64(42), uint64(57)).Return(int64(3), nil)
	app.HandleCreateOrder(w, r)

	require.Equal(t, http.StatusCreated, w.Result().StatusCode)

	res := make(map[string]interface{})
	json.NewDecoder(w.Result().Body).Decode(&res)
	assert.EqualValues(t, 3, res["order_id"])
}
