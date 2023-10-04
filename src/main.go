package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	connectionString := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		fmt.Println("Open error")
		log.Fatal(err)
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		fmt.Println("Ping error")
		log.Fatal(err)
	}

	rows, err := db.Query("select * from users")

	if err != nil {
		fmt.Println("Query error")
		log.Fatal(err)
	}

	fmt.Println(rows)
}
