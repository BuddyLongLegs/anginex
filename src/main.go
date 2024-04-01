package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/BuddyLongLegs/anginex/src/proxy"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	connStr := os.Getenv("DB_CONN_STRING")
	dbpool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	//run a simple query to check our connection
	var greeting string
	err = dbpool.QueryRow(ctx, "select 'Hello, Timescale (but concurrently)'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)

	if err := proxy.Run(); err != nil {
		log.Fatal(err)
	}

}
