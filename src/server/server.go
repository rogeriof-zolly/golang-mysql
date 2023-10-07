package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golangDB/src/database"
	"io"
	"net/http"

	"github.com/gorilla/mux"
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

func RetrieveUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	userId := params["userId"]

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when connecting to database"))
		return
	}
	defer db.Close()

	row := db.QueryRow("select * from users where id = ?", userId)

	var user user

	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("User not found!"))
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when scanning row"))
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error when converting slice to JSON"))
		return
	}
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	ID := params["userId"]

	updatedUser, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not read the request body"))
		return
	}

	var newUserData user

	if err := json.Unmarshal(updatedUser, &newUserData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not unmarshal JSON"))
		return
	}

	db, err := database.Connect()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not connect to database"))
		return
	}
	defer db.Close()

	statement, err := db.Prepare("update users set name = ?, email = ? where id = ?")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Could not create statement"))
		return
	}
	defer statement.Close()

	if _, err = statement.Exec(newUserData.Name, newUserData.Email, ID); err != nil {
		w.WriteHeader(http.StatusBadGateway)
		w.Write([]byte("Could not execute update statement"))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
