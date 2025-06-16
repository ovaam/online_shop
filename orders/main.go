package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

var (
	rdb *redis.Client
	ctx = context.Background()
)

type Order struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Status      string  `json:"status"`
}

var ordersDB = make(map[string]Order)

func main() {
	// Redis connection
	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})

	// Test Redis connection
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	r := mux.NewRouter()

	// CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// Routes
	r.HandleFunc("/orders", createOrder).Methods("POST")
	r.HandleFunc("/orders", listOrders).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrder).Methods("GET")
	r.HandleFunc("/orders/{id}/status", getOrderStatus).Methods("GET")
	r.HandleFunc("/health", healthCheck).Methods("GET")

	// Start background processor
	go processPayments()

	log.Println("Order Service started on :8080")
	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}

func createOrder(w http.ResponseWriter, r *http.Request) {
	var order Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	order.ID = "ord_" + strconv.FormatInt(time.Now().UnixNano(), 10)
	order.Status = "CREATED"
	ordersDB[order.ID] = order

	// Publish payment task
	paymentTask := map[string]interface{}{
		"order_id":    order.ID,
		"user_id":     order.UserID,
		"amount":      order.Amount,
		"description": order.Description,
	}
	if err := rdb.Publish(ctx, "payment_tasks", paymentTask).Err(); err != nil {
		http.Error(w, "Failed to process payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func listOrders(w http.ResponseWriter, r *http.Request) {
	orders := make([]Order, 0, len(ordersDB))
	for _, order := range ordersDB {
		orders = append(orders, order)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func getOrder(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order, exists := ordersDB[vars["id"]]
	if !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

func getOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order, exists := ordersDB[vars["id"]]
	if !exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": order.Status,
	})
}

func processPayments() {
	pubsub := rdb.Subscribe(ctx, "payment_results")
	defer pubsub.Close()

	ch := pubsub.Channel()
	for msg := range ch {
		var result struct {
			OrderID string `json:"order_id"`
			Status  string `json:"status"`
		}
		if err := json.Unmarshal([]byte(msg.Payload), &result); err != nil {
			log.Printf("Error decoding payment result: %v", err)
			continue
		}

		if order, exists := ordersDB[result.OrderID]; exists {
			order.Status = result.Status
			ordersDB[result.OrderID] = order
			log.Printf("Order %s updated to status %s", result.OrderID, result.Status)
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
