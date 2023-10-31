package main

import (
	"fmt"
	"log"
	"net/http"

	"amarkezic.github.com/finance-app/core"
	"amarkezic.github.com/finance-app/users"
	"github.com/gorilla/mux"
)

const port int = 9000
const url string = "localhost"

func initRouter() {
	r := mux.NewRouter()

	r.HandleFunc("/users", users.GetUsers).Methods("GET")
	r.HandleFunc("/users/{id}", users.GetUser).Methods("GET")
	r.HandleFunc("/users", users.CreateUser).Methods("POST")
	r.HandleFunc("/users/{id}", users.UpdateUser).Methods("PUT")
	r.HandleFunc("/users/{id}", users.DeleteUser).Methods("DELETE")

	log.Printf("Listening on %s:%d\n", url, port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", port), r)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	core.InitDB()
	core.InitValidation()
	initRouter()
}
