package model

// Sale represents the data from `sample_data_2.csv`.
type Sale struct {
	ID       int     `json:"id"`
	RefNo    string  `json:"ref_no"`
	SaleTime string  `json:"sale_time"`
	SoldTo   string  `json:"sold_to"`
	Amount   float64 `json:"amount"`
}

// Sales represents a slice of sale.
type Sales []Sale
