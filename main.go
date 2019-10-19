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
const (
	MaxSQLConnection = 140
	dsn              = "root:@tcp(127.0.0.1:3306)/konigle?charset=utf8mb4"
	dialect          = "mysql"
)

func main() {
	flag.Parse()
	fmt.Print(strings.Join(flag.Args(), "\n"))
	if *filename == "REQUIRED" {
		log.Fatal("File is required")
		return
	}

	// --------------------------------------------------------------------------
	// database connection setup
	// --------------------------------------------------------------------------
	db, err := sql.Open(dialect, dsn)
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

	if strings.Contains(*filename, ".json") {
		// process as json
		println("I think we should process this as json")
	} else {
		// we got a csv files
		processCSVFile(*filename, db)
	}
}

func processCSVFile(file string, db *sql.DB) {
	start := time.Now()

	csvfile, err := os.Open(file)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.Comma = ','

	i := 0
	ch := make(chan []string)
	var wg sync.WaitGroup

	// insert query
	sqlQuery := "INSERT INTO transactions(id, invoice_number, time, customer, amount, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?);"
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
			saveToDatabase(sqlQuery, r, db)
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

func saveToDatabase(sqlQuery string, r []string, db *sql.DB) {
	time.Sleep(time.Duration(1000+rand.Intn(8000)) * time.Millisecond)
	id := r[0]
	invoiceNumber := r[1]
	transactionTime := r[2]
	customer := r[3]
	amount := r[4]

	createdAt := time.Now()
	updatedAt := time.Now()

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

	_, err = stmt.Exec(id, invoiceNumber, transactionTime, customer, amount, createdAt, updatedAt)
	if err != nil {
		log.Fatalf("There was an error inserting: %v due to: %v", customer, err)
		return
	}

	tx.Commit()
}
