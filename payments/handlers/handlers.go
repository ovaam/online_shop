package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ovaam/online_shop/payments/models"
)

var accounts = make(map[string]models.Account)

func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if _, exists := accounts[account.UserID]; exists {
		respondWithError(w, http.StatusConflict, "Account already exists")
		return
	}

	account.ID = generateAccountID()
	accounts[account.UserID] = account

	respondWithJSON(w, http.StatusCreated, account)
}

func DepositHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	var deposit struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&deposit); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	account, exists := accounts[userID]
	if !exists {
		respondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	account.Balance += deposit.Amount
	accounts[userID] = account

	respondWithJSON(w, http.StatusOK, map[string]float64{
		"new_balance": account.Balance,
	})
}

func GetBalanceHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["id"]

	account, exists := accounts[userID]
	if !exists {
		respondWithError(w, http.StatusNotFound, "Account not found")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]float64{
		"balance": account.Balance,
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

func generateAccountID() string {
	return "acc_" + strconv.Itoa(len(accounts)+1)
}
