package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/ovaam/online_shop/orders/models"
)

var orders = make(map[string]models.Order)

func printWelcome() {
	color.Cyan(`
   ___  _  _  ____  ____  ____ 
  / __)( \/ )(  _ \( ___)(  _ \
 ( (__  )  (  )   / )__)  )   /
  \___)(_/\_)(_)\_)(____)(_)\_)
  
  Orders Service ready!`)
	color.Green("Available endpoints:")
	color.Yellow("POST   /orders - Create new order")
	color.Yellow("GET    /orders - List all orders")
	color.Yellow("GET    /orders/{id} - Get order status")
}

func main() {
	r := mux.NewRouter()

	// Регистрация обработчиков
	r.HandleFunc("/orders", createOrderHandler).Methods("POST")
	r.HandleFunc("/orders", listOrdersHandler).Methods("GET")
	r.HandleFunc("/orders/{id}", getOrderStatusHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	printWelcome()
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Улучшенный вывод ошибок
func respondWithError(w http.ResponseWriter, code int, message string) {
	color.Red("ERROR: %s", message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

// Красивый JSON вывод
func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.MarshalIndent(payload, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Обновляем обработчики
func createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Генерируем ID (в реальном приложении используйте UUID)
	order.ID = "ord_" + strconv.Itoa(len(orders)+1)
	order.Status = models.StatusNew

	orders[order.ID] = order
	color.Green("Created new order: %s", order.ID)

	// TODO: Добавить асинхронную обработку оплаты
	respondWithJSON(w, http.StatusCreated, order)
}

func listOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orderList []models.Order
	for _, order := range orders {
		orderList = append(orderList, order)
	}
	respondWithJSON(w, http.StatusOK, orderList)
}

func getOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderID := vars["id"]

	order, exists := orders[orderID]
	if !exists {
		respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"status": order.Status,
		"amount": order.Amount,
	})
}
