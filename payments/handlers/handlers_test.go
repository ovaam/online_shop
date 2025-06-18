package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ovaam/online_shop/payments/models"
)

func TestCreateAccountHandler(t *testing.T) {
	req, err := http.NewRequest("POST", "/accounts", strings.NewReader(`{"user_id":"test1"}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateAccountHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}
}

func TestDepositHandler(t *testing.T) {
	// First create an account
	account := models.Account{
		ID:      "test_account",
		UserID:  "test_user",
		Balance: 0,
	}
	accounts[account.UserID] = account

	req, err := http.NewRequest("POST", "/accounts/test_user/deposit", strings.NewReader(`{"amount":100}`))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DepositHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestGetBalanceHandler(t *testing.T) {
	// First create an account
	account := models.Account{
		ID:      "test_account",
		UserID:  "test_user",
		Balance: 100,
	}
	accounts[account.UserID] = account

	req, err := http.NewRequest("GET", "/accounts/test_user/balance", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetBalanceHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}
