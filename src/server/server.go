package server

import (
	"encoding/json"
	"fmt"
	"golangDB/src/database"
	"io"
	"net/http"
)

type user struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	newUser, err := io.ReadAll(r.Body)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("It was not possible to read request body"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error connecting to the database"))
		return
	}
	defer db.Close()

	var user user

	if err = json.Unmarshal(newUser, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when converting user to JSON"))
		return
	}

	statement, err := db.Prepare("insert into users (name, email) values (?, ?)")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when creating statement"))
		return
	}
	defer statement.Close()

	insertion, err := statement.Exec(user.Name, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when inserting new user"))
		return
	}

	insertedId, err := insertion.LastInsertId()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when retrieving ID"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created successfully! Id: %d", insertedId)))
}

func RetrieveUsers(w http.ResponseWriter, r *http.Request) {

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when connecting to database"))
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from users")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when retrieving data from DB"))
		return
	}
	defer rows.Close()

	var users []user

	for rows.Next() {
		var user user

		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("Error when scanning row"))
			return
		}

		users = append(users, user)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(users); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when converting slice to JSON"))
		return
	}
}
