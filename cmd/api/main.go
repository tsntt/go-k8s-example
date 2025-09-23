package main

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type User struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

var (
	httpRequestsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Total HTTP requests.",
	})
	// atomic flag for readiness probe
	isReady atomic.Bool
	db      *sqlx.DB
)

func initDB() error {
	connStr := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	var err error
	db, err = sqlx.Connect("postgres", connStr)
	if err != nil {
		return fmt.Errorf("Failed to connect to db: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("failed to ping db: %w", err)
	}

	schema := `CREATE TABLE IF NOT EXISTS users (
		id VARCHAR(255) PRIMARY KEY,
		name VARCHAR(255)
	);`
	if _, err := db.Exec(schema); err != nil {
		return fmt.Errorf("Failed to create table: %w", err)
	}
	slog.Info("db connection established sucessfully")
	return nil
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
		var users []User
		err := db.Select(&users, "SELECT * FROM users")
		if err != nil {
			http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)

	case "POST":
		var newUser User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_, err := db.NamedExec("INSERT INTO users (id, name) VALUES (:id, :name)", newUser)
		if err != nil {
			http.Error(w, "Failed inserting user", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(fmt.Sprintf("User %s successfully created!", newUser.Name)))

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

	slog.Info("Server running", "port", 8000)
	err := http.ListenAndServe(":8000", mux)
	if err != nil {
		slog.Error("Failed to start the server", "error", err)
	}
}
