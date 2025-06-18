package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ovaam/online_shop/orders/models"
)

var orders = make(map[string]models.Order)

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	var order models.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order.ID = generateOrderID()
	order.Status = models.StatusNew
	orders[order.ID] = order

	respondWithJSON(w, http.StatusCreated, order)
}

func ListOrdersHandler(w http.ResponseWriter, r *http.Request) {
	var orderList []models.Order
	for _, order := range orders {
		orderList = append(orderList, order)
	}
	respondWithJSON(w, http.StatusOK, orderList)
}

func GetOrderStatusHandler(w http.ResponseWriter, r *http.Request) {
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

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func generateOrderID() string {
	return "ord_" + strconv.Itoa(len(orders)+1)
}
