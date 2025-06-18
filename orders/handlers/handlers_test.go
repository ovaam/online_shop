package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ovaam/online_shop/orders/models"
)

func TestCreateOrderHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/orders", strings.NewReader(`{"user_id":"test1","amount":100,"description":"Test order"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateOrderHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestListOrdersHandler(t *testing.T) {
	req, err := http.NewRequest("GET", "/orders", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListOrdersHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetOrderStatusHandler(t *testing.T) {
	// First create an order
	order := models.Order{
		ID:          "test_order",
		UserID:      "test_user",
		Amount:      100,
		Description: "Test order",
		Status:      models.StatusNew,
	}
	orders[order.ID] = order

	req, err := http.NewRequest("GET", "/orders/test_order", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetOrderStatusHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
