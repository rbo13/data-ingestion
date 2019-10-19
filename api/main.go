package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rbo13/data-ingestion/api/handler"
	"github.com/rbo13/data-ingestion/api/router"
	"github.com/rbo13/data-ingestion/api/server"

	_ "github.com/go-sql-driver/mysql"
)

const (
	// MethodGet ...
	MethodGet = "GET"
	// MethodPost ...
	MethodPost = "POST"
	// MethodDelete ...
	MethodDelete = "DELETE"
	// MethodPut ...
	MethodPut = "PUT"
	// MethodPatch ...
	MethodPatch = "PATCH"
	//MaxSQLConnection ...
	MaxSQLConnection = 140

	dsn     = "root:@tcp(127.0.0.1:3306)/konigle?charset=utf8mb4"
	dialect = "mysql"
)

func main() {
	router := router.New()

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

	// routes
	router.HandleFunc(MethodGet, "/api/v1/transactions", handler.TransactionHandler(db))
	router.HandleFunc(MethodGet, "/api/v1/transactions/:id", handler.SingleTransactionHandler(db))

	srv := server.New(":1337", router)
	go func() {
		srv.Start()
	}()

	gracefulShutdown(srv.HTTPServer)
}

func gracefulShutdown(srv *http.Server) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
