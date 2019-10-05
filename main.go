package main

import (
	"database/sql"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var filename = flag.String("f", "REQUIRED", "source CSV file")
var numChannels = flag.Int("c", 4, "num of parallel channels")

// MaxSQLConnection represents the maximum number of SQL Connection.
// MySQL allows only 150 connections per session.
const MaxSQLConnection = 140

func main() {
	start := time.Now()
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), "\n"))
	if *filename == "REQUIRED" {
		return
	}

	csvfile, err := os.Open(*filename)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.Comma = ','

	// --------------------------------------------------------------------------
	// database connection setup
	// --------------------------------------------------------------------------
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/konigle?charset=utf8mb4")
	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// check database connection
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	// set max idle connections
	db.SetMaxIdleConns(MaxSQLConnection)
	defer db.Close()
	println("Successfully Connected to Database!")

	i := 0
	ch := make(chan []string)
	var wg sync.WaitGroup

	// insert query
	sqlQuery := "INSERT INTO transactions(id, invoice_number, time, customer, amount) VALUES (?, ?, ?, ?, ?);"

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return
		}
		i++

		wg.Add(1)
		go func(r []string, i int) {
			defer wg.Done()
			processData(sqlQuery, r, db)
			ch <- r
		}(record, i)

		// fmt.Printf("\rgo %d", i)
	}

	// closer
	go func() {
		wg.Wait()
		close(ch)
	}()

	// print channel results (necessary to prevent exit programm before)
	j := 0
	for range ch {
		j++
		fmt.Printf("\r\t\t\t\t | done %d \n", j)
	}

	fmt.Printf("\n%2fs", time.Since(start).Seconds())
}

func processData(sqlQuery string, r []string, db *sql.DB) {
	time.Sleep(time.Duration(1000+rand.Intn(8000)) * time.Millisecond)
	id := r[0]
	invoiceNumber := r[1]
	transactionTime := r[2]
	customer := r[3]
	amount := r[4]

	log.Printf("\n ID: %v -> INVOICE NUMBER: %v -> TRANSACTION TIME: %v -> CUSTOMER: %v -> AMOUNT: %v \n", id, invoiceNumber, transactionTime, customer, amount)

	// save to database
	if len(r) <= 0 {
		return
	}

	tx, err := db.Begin()
	if err != nil {
		return
	}

	stmt, err := tx.Prepare(sqlQuery)
	if err != nil {
		log.Println(err)
		return
	}

	defer stmt.Close()

	_, err = stmt.Exec(id, invoiceNumber, transactionTime, customer, amount)
	log.Println(err)
	if err != nil {
		return
	}

	// fmt.Printf("\r\t\t| proc %d", i)
}
