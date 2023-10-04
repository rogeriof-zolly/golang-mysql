package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type user struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func CreateUser(w http.ResponseWriter, r http.Request) {

	newUser, err := io.ReadAll(r.Body)

	if err != nil {
		log.Fatal(err)
		return
	}

	var user user

	if err = json.Unmarshal(newUser, &user); err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(user)

}
