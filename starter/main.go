package main

import (
	"context"
	"log"
	"my_module/expense"
	"my_module/types"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.temporal.io/sdk/client"
)

// Global variables
// - mutex: a mutex to protect the expenseData map
// - expenseData: a map that stores the expense reports
var (
	HTTPPort = "8097"
	// mutex    sync.Mutex
	temporal client.Client
)

func main() {
	var err error
	temporal, err = client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("unable to create Temporal client", err)
	}
	log.Println("Temporal client connected")

	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/expense", CreateExpenseHandler)
	router.GET("/expense/:userid", QueryExpenseHandler)

	router.Run("localhost:" + HTTPPort)
}

// CreateExpenseHandler is the handler for the /create endpoint. It creates a new expense report by starting the CreateExpense workflow using temporal. Once the expense has been
// created, it return the expenseID to the client in JSON format. This handler uses gin to communicate with the client.
func CreateExpenseHandler(c *gin.Context) {
	var newExpense types.ExpenseReport
	c.BindJSON(&newExpense)

	newExpense.ExpenseID = time.Now().Format(time.RFC3339)

	workflowOptions := client.StartWorkflowOptions{
		ID:        "createExpense_" + newExpense.ExpenseID,
		TaskQueue: "expense",
	}

	we, err := temporal.ExecuteWorkflow(context.Background(), workflowOptions, expense.CreateExpenseWorkflow, &newExpense)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
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
