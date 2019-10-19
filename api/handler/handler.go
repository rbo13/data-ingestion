package handler

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/rbo13/data-ingestion/api/model"
	"github.com/rbo13/data-ingestion/api/router"
)

// TransactionHandler handles transaction...
func TransactionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		transactions := getAllTransactions(db)
		json.NewEncoder(w).Encode(transactions)
	}
}

// SingleTransactionHandler handles the querying of single transaction using a given transaction id.
func SingleTransactionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		transactionID, err := strconv.Atoi(router.Param(r.Context(), "id"))
		if err != nil {
			log.Println(err)
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")

		transaction := getTransaction(db, transactionID)
		json.NewEncoder(w).Encode(transaction)
	}
}

func getTransaction(db *sql.DB, id int) model.Transaction {

	var transaction model.Transaction

	tx, err := db.Begin()
	if err != nil {
		return transaction
	}

	sqlQuery := "SELECT * FROM transactions WHERE id = ?;"

	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return transaction
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		log.Fatalf("There was an error inserting: %v due to: %v", transaction, err)
		return transaction
	}

	err = stmt.QueryRow(id).Scan(
		&transaction.ID,
		&transaction.Time,
		&transaction.InvoiceNumber,
		&transaction.Customer,
		&transaction.Amount,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		log.Fatalf("Error Due to: %v", err)
		return transaction
	}

	// log.Printf("Transaction: %v", transaction)
	return transaction
}

func getAllTransactions(db *sql.DB) model.Transactions {

	var transactions model.Transactions
	var transaction model.Transaction

	tx, err := db.Begin()
	if err != nil {
		return transactions
	}

	sqlQuery := "SELECT * FROM transactions;"

	results, err := tx.Query(sqlQuery)
	if err != nil {
		log.Println(err)
		return transactions
	}

	defer results.Close()

	for results.Next() {
		err = results.Scan(
			&transaction.ID,
			&transaction.Time,
			&transaction.InvoiceNumber,
			&transaction.Customer,
			&transaction.Amount,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		transactions = append(transactions, transaction)
	}

	// log.Printf("Transaction: %v", transaction)
	return transactions
}
