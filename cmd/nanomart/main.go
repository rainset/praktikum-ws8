package main

import (
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/kilchik/nanomart/pkg/api"
	"github.com/kilchik/nanomart/pkg/storage"
)

const addr = "0.0.0.0:5301"

func main() {
	db := sqlx.MustOpen("sqlite3", "nanomart.db")

	store := storage.New(db)
	app := api.New(store)

	http.HandleFunc("/api/v1/order", app.HandleCreateOrder)

	log.Printf("listening %v", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("listen %s: %v", addr, err)
	}
}
