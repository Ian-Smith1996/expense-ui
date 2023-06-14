package main

import (
	"log"
	"my_module/expense"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// Main function that creates a temporal worker and registers it to run the CreateExpenseWorkflow workflow and the CreateExpenseActivity activity.
func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "expense", worker.Options{})

	w.RegisterWorkflow(expense.CreateExpenseWorkflow)
	w.RegisterActivity(expense.CreateExpenseActivity)
	// w.RegisterActivity(expense.WaitForDecisionActivity)
	// w.RegisterActivity(expense.PaymentActivity)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
