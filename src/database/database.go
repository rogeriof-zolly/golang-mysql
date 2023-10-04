package database

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() (*sql.DB, error) {
	connectionString := "golang:golang@/devbook?charset=utf8&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		log.Print("Open database error")
		return nil, err
	}

	err = db.Ping()

	if err != nil {
		log.Print("Ping error")
		return nil, err
	}

	return db, nil
}
