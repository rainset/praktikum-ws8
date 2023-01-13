package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/kilchik/nanomart/internal/pkg/metrix"
	tracerpkg "github.com/kilchik/nanomart/internal/pkg/tracer"
	"github.com/kilchik/nanomart/pkg/api"
	"github.com/kilchik/nanomart/pkg/storage"
)

const (
	addrApp     = "0.0.0.0:5301"
	addrMetrics = "0.0.0.0:5302"
)

const jaegerCollectorURL = "http://localhost:14268/api/traces"

func main() {
	// Initialize metrics server
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	metricsSrv := http.Server{
		Addr:    addrMetrics,
		Handler: router,
	}
	go metricsSrv.ListenAndServe()

	// Initialize tracer
	tracer, err := tracerpkg.Build(jaegerCollectorURL)
	if err != nil {
		log.Fatalf("build tracer: %v", err)
	}
	defer tracer.Shutdown()

	_ = metrix.New()

	db := sqlx.MustOpen("sqlite3", "nanomart.db")

	store := storage.New(db)
	app := api.New(store)

	http.HandleFunc("/api/v1/order", tracerpkg.Middleware(app.HandleCreateOrder))

	log.Printf("listening %v", addrApp)
	if err := http.ListenAndServe(addrApp, nil); err != nil {
		log.Fatalf("listen %s: %v", addrApp, err)
	}
}
