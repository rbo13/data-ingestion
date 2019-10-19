package model

// Transaction represents the transactions table.
type Transaction struct {
	ID            int     `json:"id"`
	Time          string  `json:"time"`
	InvoiceNumber string  `json:"invoice_number"`
	Customer      string  `json:"customer"`
	Amount        float32 `json:"amount"`
	CreatedAt     []uint8 `json:"created_at"`
	UpdatedAt     []uint8 `json:"updated_at"`
}

// Transactions is a slice of Transaction.
type Transactions []Transaction
