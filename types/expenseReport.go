package types

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
