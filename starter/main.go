package main

import (
	"context"
	"log"
	"my_module/expense"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

// ExpenseReport is the struct that represents an expense report. It contains the following fields:
// - ExpenseID: a unique identifier for the expense report
// - Amount: the amount of the expense
// - Date: the date of the expense
// - State: the state of the expense report. It can be "created", "approved", or "rejected".
type ExpenseReport struct {
	ExpenseID string  `json:"expenseID"`
	Amount    float64 `json:"amount,omitempty"`
	Date      string  `json:"date,omitempty"`
	State     string  `json:"state,omitempty"`
}

// Global variables
// - mutex: a mutex to protect the expenseData map
// - expenseData: a map that stores the expense reports
var (
	HTTPPort    = "8097"
	mutex       sync.Mutex
	expenseData = make(map[string]ExpenseReport)
	temporal    client.Client
)

func main() {
	var err error
	temporal, err = client.NewClient(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	log.Println("Temporal client connected")

	router := gin.Default()
	router.POST("/expense", CreateExpenseHandler)
	router.GET("/expense/:userid", QueryExpenseHandler)

	router.Run("localhost:" + HTTPPort)
}

// CreateExpenseHandler is the handler for the /create endpoint. It creates a new expense report by starting the CreateExpense workflow using temporal. Once the expense has been
// created, it return the expenseID to the client in JSON format. This handler uses gin to communicate with the client.
func CreateExpenseHandler(c *gin.Context) {
	var newExpense ExpenseReport

	err := c.BindJSON(&newExpense)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newExpense.ExpenseID = time.Now().Format(time.RFC3339)
	newExpense.State = "created"

	mutex.Lock()
	expenseData[newExpense.ExpenseID] = newExpense
	mutex.Unlock()

	workflowOptions := client.StartWorkflowOptions{
		ID:        "createExpense_" + newExpense.ExpenseID,
		TaskQueue: "expense",
	}

	we, err := temporal.ExecuteWorkflow(context.Background(), workflowOptions, expense.ExpenseWorkflow, newExpense)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	log.Println("Started workflow", "WorkflowID", we.GetID(), "RunID", we.GetRunID())

	c.JSON(http.StatusOK, gin.H{"expenseID": newExpense.ExpenseID})

}

// QueryExpenseHandler is the handler for the /query endpoint. It uses gin to extract the id from the URL then queries the running workflow with that ID for it's state.
// It then sends the state back to the client as a JSON.
func QueryExpenseHandler(c *gin.Context) {
	expenseID := c.Param("userid")

	queryResult, err := temporal.QueryWorkflow(context.Background(), expenseID, "ExpenseState", "state")
	if err != nil {
		log.Fatalln("Unable to query workflow", err)
	}

	var state string
	err = queryResult.Get(&state)
	if err != nil {
		log.Fatalln("Unable to get query result", err)
	}

	c.JSON(http.StatusOK, gin.H{"state": state})
}
