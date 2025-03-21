package metrick

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net/http"
)

var (
	RequestDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_duration_seconds",
		Help:    "Время обработки HTTP запросов",
		Buckets: []float64{0.1, 0.25, 0.5, 1, 2.5},
	}, []string{"handler", "method"})

	RequestCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Общее количество HTTP запросов",
	}, []string{"handler", "method", "status"})
)

// Помощник для записи статуса ответа
type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (rec *StatusRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.ResponseWriter.WriteHeader(code)
}
