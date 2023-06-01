package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"my_module/expense"
	"net/http"
	"sync"
	"time"

	"go.temporal.io/sdk/client"
)

// ExpenseReport is the struct that represents an expense report. It contains the following fields:
// - ExpenseID: a unique identifier for the expense report
// - Amount: the amount of the expense
// - Date: the date of the expense
// - State: the state of the expense report. It can be "created", "approved", or "rejected".
type ExpenseReport struct {
	ExpenseID string  `json:"expenseID"`
	Amount    float64 `json:"amount"`
	Date      string  `json:"date"`
	State     string  `json:"state,omitempty"`
}

// Global variables
// - mutex: a mutex to protect the expenseData map
// - expenseData: a map that stores the expense reports
var (
	mutex       sync.Mutex
	expenseData = make(map[string]ExpenseReport)
)

// CreateExpenseReport is the handler for the /create endpoint. It should receive a POST request with a JSON body that contains
// an ExpenseReport object. It will store the ExpenseReport object in expenseData and return the ExpenseID in the response if it
// is successful. Otherwise, it will send an error response with a 500 error code. If the ExpenseReport is saved successfully,
// it will start a workflow: CreateExpenseWorkflow using temporal.
func CreateExpenseReport(w http.ResponseWriter, r *http.Request) {
	// Create a temporal client and defer closing it
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
		return
	}
	defer c.Close()

	// Handling CORS here
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	var newExpense ExpenseReport

	// Decode the request body into an ExpenseReport object
	err = json.NewDecoder(r.Body).Decode(&newExpense)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Set the ExpenseID and State fields
	newExpense.ExpenseID = time.Now().Format(time.RFC3339)
	newExpense.State = "created"

	// Store the ExpenseReport object in expenseData while being async safe
	mutex.Lock()
	expenseData[newExpense.ExpenseID] = newExpense
	fmt.Println("Received and created expense with ID " + expenseData[newExpense.ExpenseID].ExpenseID)
	mutex.Unlock()

	response := "{\"expenseID\":\"" + newExpense.ExpenseID + "\"}"

	// Start a workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        "expense_" + newExpense.ExpenseID,
		TaskQueue: "expense",
	}
	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, expense.SampleExpenseWorkflow, newExpense)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Start a HTTP server the listens on port 8098 and creates a new expense report when it receives a POST request to `/create
	http.HandleFunc("/create", CreateExpenseReport)
	fmt.Println("Listening...")
	http.ListenAndServe(":8098", nil)
}
