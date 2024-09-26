package main

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "user=postgres dbname=medods password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Connected!")
}
