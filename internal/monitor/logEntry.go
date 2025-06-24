package monitor

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type LogEntry struct {
	Timestamp string        `json:"timestamp"`
	RemoteIP  string        `json:"remote_ip"`
	Method    string        `json:"method"`
	Path      string        `json:"path"`
	Target    string        `json:"target,omitempty"`
	Duration  time.Duration `json:"duration"`
	Status    int           `json:"status"`
	UserAgent string        `json:"user_agent"`
	Referer   string        `json:"referer"`
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (w *statusWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

var (
	httpRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "apigw_http_requests_total",
			Help: "Total de peticiones HTTP",
		},
		[]string{"path", "method", "status"},
	)
	httpDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "apigw_http_request_duration_seconds",
			Help:    "Duración de las peticiones HTTP",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequests, httpDuration)
}

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: 200}
		next.ServeHTTP(sw, r)
		target, _ := r.Context().Value("target").(string)
		entry := LogEntry{
			Timestamp: time.Now().Format(time.RFC3339),
			RemoteIP:  r.RemoteAddr,
			Method:    r.Method,
			Path:      r.URL.Path,
			Target:    target,
			Status:    sw.status,
			Duration:  time.Since(start),
			UserAgent: r.UserAgent(),
			Referer:   r.Referer(),
		}
		enc := json.NewEncoder(os.Stdout)
		enc.Encode(entry)
	})
}

func PrometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sw := &statusWriter{ResponseWriter: w, status: 200}
		timer := prometheus.NewTimer(httpDuration.WithLabelValues(r.URL.Path, r.Method))
		next.ServeHTTP(sw, r)
		timer.ObserveDuration()
		httpRequests.WithLabelValues(r.URL.Path, r.Method, http.StatusText(sw.status)).Inc()
	})
}

func PrometheusHandler() http.Handler {
	return promhttp.Handler()
}
