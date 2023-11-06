package main

import (
	"fmt"
	"log"
	"net/http"

	"amarkezic.github.com/finance-app/company"
	"amarkezic.github.com/finance-app/core"
	"amarkezic.github.com/finance-app/users"
	"github.com/gorilla/mux"
)

const port int = 9000
const url string = "localhost"

func initRouter() {
	r := mux.NewRouter()

	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(core.AuthMiddleware)
	userRoutes.HandleFunc("", users.GetUsers).Methods("GET")
	userRoutes.HandleFunc("/{id}", users.GetUser).Methods("GET")
	userRoutes.HandleFunc("", users.CreateUser).Methods("POST")
	userRoutes.HandleFunc("/{id}", users.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id}", users.DeleteUser).Methods("DELETE")

	companyRoutes := r.PathPrefix("/companies").Subrouter()
	companyRoutes.Use(core.AuthMiddleware)
	companyRoutes.HandleFunc("", company.GetCompanies).Methods("GET")
	companyRoutes.HandleFunc("/{id}", company.GetCompany).Methods("GET")
	companyRoutes.HandleFunc("", company.CreateCompany).Methods("POST")
	companyRoutes.HandleFunc("/{id}", company.UpdateCompany).Methods("PUT")
	companyRoutes.HandleFunc("/{id}", company.DeleteCompany).Methods("DELETE")

	r.HandleFunc("/auth", users.Login).Methods("POST")

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
