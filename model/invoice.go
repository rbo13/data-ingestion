package model

// Invoice represents the invoice from `sample_data_1.csv` file.
type Invoice struct {
	ID        int    `json:"id"`
	InvoiceNo string `json:"invoice_no"`
	Time      string `json:"time"`
	Customer  string `json:"customer"`
	Amount    string `json:"amount"`
}

// Invoices custom type that represents a slice of Invoice.
type Invoices []Invoice
