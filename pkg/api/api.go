package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
)

type Storage interface {
	InsertOrder(ctx context.Context, userID uint64, total uint64) (int64, error)
}

type Item struct {
	Name  string `json:"name"`
	Price uint64 `json:"price"`
}

type CreateOrderRequest struct {
	UserID *uint64 `json:"user_id"`
	Items  []Item  `json:"items"`
}

type App struct {
	store Storage
}

func New(store Storage) *App {
	return &App{store: store}
}

func (a *App) HandleCreateOrder(w http.ResponseWriter, r *http.Request) {
	// Decode request
	req := &CreateOrderRequest{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		log.Printf("[E] decode request: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Validate request
	if req.UserID == nil || len(req.Items) == 0 {
		log.Printf("[E] invalid request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Calculate data for DB
	var total uint64
	for _, i := range req.Items {
		total += i.Price
	}

	// Save data to DB
	orderID, err := a.store.InsertOrder(r.Context(), *req.UserID, total)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Reply
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{"order_id": orderID})
}
