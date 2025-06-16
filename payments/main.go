package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
	"github.com/ovaam/online_shop/payments/models"
)

var accounts = make(map[string]models.Account)

func printWelcome() {
	color.Cyan(`
  ____  _____  __  __  ____  ____  ____ 
 ( ___)(  _  )(  \/  )( ___)(  _ \( ___)
  )__)  )(_)(  )    (  )__)  )   / )__) 
 (____)(_____)(_/\/\_)(____)(_)\_)(____)
  
  Payments Service ready!`)
	color.Green("Available endpoints:")
	color.Yellow("POST   /accounts - Create new account")
	color.Yellow("POST   /accounts/{id}/deposit - Deposit funds")
	color.Yellow("GET    /accounts/{id}/balance - Check balance")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/accounts", createAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}/deposit", depositHandler).Methods("POST")
	r.HandleFunc("/accounts/{id}/balance", getBalanceHandler).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	printWelcome()
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Аналогичные улучшения для ошибок и JSON вывода
func respondWithError(w http.ResponseWriter, code int, message string) {
	color.Red("ERROR: %s", message)
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.MarshalIndent(payload, "", "  ")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// Обновляем обработчики
func createAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account models.Account
	if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Проверяем, есть ли уже аккаунт
	if _, exists := accounts[account.UserID]; exists {
		respondWithError(w, http.StatusConflict, "Account already exists")
		return
	}

	account.ID = "acc_" + strconv.Itoa(len(accounts)+1)
	accounts[account.UserID] = account

	color.Green("Created new account for user: %s", account.UserID)
	respondWithJSON(w, http.StatusCreated, account)
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
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

	color.Blue("Deposited %.2f to account %s", deposit.Amount, userID)
	respondWithJSON(w, http.StatusOK, map[string]float64{
		"new_balance": account.Balance,
	})
}

func getBalanceHandler(w http.ResponseWriter, r *http.Request) {
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
