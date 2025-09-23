package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests.",
	})
	// atomic flag for readiness probe
	isReady atomic.Bool
)

var users = []User{
	{ID: "1", Name: "Alice"},
	{ID: "2", Name: "Bob"},
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	if isReady.Load() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
		return
	}
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("not ready"))
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("Endpoint /users accessed", "method", r.Method)
	httpRequestsTotal.Inc()

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readiness", readinessHandler)
	mux.HandleFunc("/users", usersHandler)

	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("Fake initialization of the application...")
	go func() {
		time.Sleep(5 * time.Second)
		isReady.Store(true)
		slog.Info("Application is ready to receive traffic.")
	}()

	slog.Info("Server running", "port", 8080)
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		slog.Error("Failed to start the server", "error", err)
	}
}
