package main

import (
	"encoding/json"
	"fmt"
	"my_module/types"
	"net/http"
	"sync"
	"time"
)

type ExpenseResponse struct {
	ExpenseID string `json:"expenseID"`
	State     string `json:"state"`
}

type ExpenseQuery struct {
	ExpenseID string `json:"expenseID"`
}

var (
	mutex       sync.Mutex
	expenseData = make(map[string]types.ExpenseReport)
)

func CreateExpenseReport(w http.ResponseWriter, r *http.Request) {
	// Handling CORS here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var newExpense types.ExpenseReport

	err := json.NewDecoder(r.Body).Decode(&newExpense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newExpense.ExpenseID = time.Now().Format(time.RFC3339)
	newExpense.State = "created"

	mutex.Lock()
	expenseData[newExpense.ExpenseID] = newExpense
	fmt.Println("Received and created expense with ID " + expenseData[newExpense.ExpenseID].ExpenseID)
	mutex.Unlock()

	response := ExpenseResponse{
		ExpenseID: newExpense.ExpenseID,
		State:     newExpense.State,
	}

	json.NewEncoder(w).Encode(response)
}
func QueryExpense(w http.ResponseWriter, r *http.Request) {
	// Handling CORS here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var expenseID types.ExpenseReport
	fmt.Println(r.Body)
	err := json.NewDecoder(r.Body).Decode(&expenseID)
	if err != nil {
		fmt.Println("error parsing json")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	expense, ok := expenseData[expenseID.ExpenseID]

	if !ok {
		http.Error(w, "Expense not found for id: "+expenseID.ExpenseID, http.StatusNotFound)
		return
	}
	expense.State = "TODO"
	response := ExpenseResponse{
		State: expense.State,
	}

	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/create", CreateExpenseReport)
	http.HandleFunc("/query", QueryExpense)
	fmt.Println("Listening...")
	http.ListenAndServe(":8098", nil)
}
