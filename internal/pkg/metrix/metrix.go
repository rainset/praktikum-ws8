package metrix

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	NameCounterCreateResults = "create_order_results"
	NameSummaryCreateLatency = "create_order_latency"
)

type Metrix struct {
	counterCreateResults *prometheus.CounterVec
	summaryCreateLatency *prometheus.SummaryVec
}

func New() *Metrix {
	return &Metrix{
		counterCreateResults: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: NameCounterCreateResults,
				Help: "HTTP status codes of create-order method",
			}, []string{"code"},
		),
		summaryCreateLatency: promauto.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       NameSummaryCreateLatency,
				Help:       "Duration of create-order method",
				Objectives: map[float64]float64{0.5: 0.5, 0.9: 0.9, 1: 1},
				MaxAge:     30 * time.Second,
			}, []string{"method"},
		),
	}
}

func (m *Metrix) IncResultsCounter(code int) {
	m.counterCreateResults.With(map[string]string{"code": strconv.Itoa(code)}).Inc()
}

func (m *Metrix) ObserveLatency(start time.Time) {
	dur := float64(time.Since(start).Milliseconds())
	m.summaryCreateLatency.With(map[string]string{"method": "create"}).Observe(dur)
}
