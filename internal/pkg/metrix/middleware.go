package metrix

import (
	"net/http"
	"time"
)

func (m *Metrix) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w}

		start := time.Now()
		next(rec, r)

		m.ObserveLatency(start)
		m.IncResultsCounter(rec.code)
	}
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.code = code
	sr.ResponseWriter.WriteHeader(code)
}
