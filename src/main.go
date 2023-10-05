package main

import (
	"golangDB/src/server"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/users", server.CreateUser).Methods(http.MethodPost)
	router.HandleFunc("/users", server.RetrieveUsers).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":5000", router))
}
