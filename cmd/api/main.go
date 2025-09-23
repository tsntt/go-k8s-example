package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Estrutura de dados para o usuário
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

var (
	// Métrica para monitorar o total de requisições HTTP
	httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total de requisições HTTP.",
	})
	// Flag atômica para o readiness probe
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
	slog.Info("Endpoint /users acessado", "method", r.Method)
	httpRequestsTotal.Inc() // Incrementa a métrica de requisições

	switch r.Method {
	case "GET":
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func main() {
	// Cria um logger estruturado em JSON
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/readiness", readinessHandler)
	mux.HandleFunc("/users", usersHandler)

	// Expondo as métricas do Prometheus
	mux.Handle("/metrics", promhttp.Handler())

	slog.Info("Simulando inicialização da aplicação...")
	go func() {
		time.Sleep(5 * time.Second) // Simula uma tarefa demorada
		isReady.Store(true)
		slog.Info("Aplicação está pronta para receber tráfego.")
	}()

	slog.Info("Servidor iniciado", "port", 8080)
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		slog.Error("Falha ao iniciar o servidor", "error", err)
	}
}