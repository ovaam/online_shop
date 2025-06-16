package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	rdb      *redis.Client
	ctx      = context.Background()
	accounts = make(map[string]Account)
	mu       sync.Mutex
)

type Account struct {
	ID      string  `json:"id"`
	UserID  string  `json:"user_id"`
	Balance float64 `json:"balance"`
}

func main() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	r.HandleFunc("/accounts", createAccount).Methods("POST")
	r.HandleFunc("/accounts/{user_id}/deposit", deposit).Methods("POST")
	r.HandleFunc("/accounts/{user_id}/balance", getBalance).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Start payment processor
	go processPaymentTasks()

	log.Println("Payment Service started on :8081")
	log.Fatal(http.ListenAndServe(":8081", c.Handler(r)))
}

func createAccount(w http.ResponseWriter, r *http.Request) {
	var acc Account
	if err := json.NewDecoder(r.Body).Decode(&acc); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	if _, exists := accounts[acc.UserID]; exists {
		http.Error(w, "Account already exists", http.StatusConflict)
		return
	}

	acc.ID = "acc_" + acc.UserID
	accounts[acc.UserID] = acc

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func deposit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	mu.Lock()
	defer mu.Unlock()

	acc, exists := accounts[userID]
	if !exists {
		http.NotFound(w, r)
		return
	}

	acc.Balance += req.Amount
	accounts[userID] = acc

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(acc)
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["user_id"]

	mu.Lock()
	acc, exists := accounts[userID]
	mu.Unlock()

	if !exists {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]float64{
		"balance": acc.Balance,
	})
}

func processPaymentTasks() {
	pubsub := rdb.Subscribe(ctx, "payment_tasks")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var task struct {
			OrderID string  `json:"order_id"`
			UserID  string  `json:"user_id"`
			Amount  float64 `json:"amount"`
		}
		if err := json.Unmarshal([]byte(msg.Payload), &task); err != nil {
			log.Printf("Error decoding payment task: %v", err)
			continue
		}

		mu.Lock()
		acc, exists := accounts[task.UserID]
		mu.Unlock()

		result := make(map[string]interface{})
		result["order_id"] = task.OrderID

		if !exists {
			result["status"] = "FAILED"
			result["reason"] = "Account not found"
		} else if acc.Balance < task.Amount {
			result["status"] = "FAILED"
			result["reason"] = "Insufficient funds"
		} else {
			mu.Lock()
			acc.Balance -= task.Amount
			accounts[task.UserID] = acc
			mu.Unlock()
			result["status"] = "PAID"
		}

		if err := rdb.Publish(ctx, "payment_results", result).Err(); err != nil {
			log.Printf("Failed to publish payment result: %v", err)
		}
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	if err := rdb.Ping(ctx).Err(); err != nil {
		http.Error(w, "Redis unavailable", http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
