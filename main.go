package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"amarkezic.github.com/finance-app/company"
	"amarkezic.github.com/finance-app/core"
	"amarkezic.github.com/finance-app/projects"
	"amarkezic.github.com/finance-app/records"
	"amarkezic.github.com/finance-app/users"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func initRouter() *mux.Router {
	r := mux.NewRouter()

	userRoutes := r.PathPrefix("/users").Subrouter()
	userRoutes.Use(core.AuthMiddleware)
	userRoutes.Use(core.AuthorizationMiddleware)
	userRoutes.HandleFunc("", users.GetUsers).Methods("GET")
	userRoutes.HandleFunc("/{id}", users.GetUser).Methods("GET")
	userRoutes.HandleFunc("/{id}", users.UpdateUser).Methods("PUT")
	userRoutes.HandleFunc("/{id}", users.DeleteUser).Methods("DELETE")

	companyRoutes := r.PathPrefix("/companies").Subrouter()
	companyRoutes.Use(core.AuthMiddleware)
	companyRoutes.Use(core.AuthorizationMiddleware)
	companyRoutes.HandleFunc("", company.GetCompanies).Methods("GET")
	companyRoutes.HandleFunc("/{id}", company.GetCompany).Methods("GET")
	companyRoutes.HandleFunc("", company.CreateCompany).Methods("POST")
	companyRoutes.HandleFunc("/{id}", company.UpdateCompany).Methods("PUT")
	companyRoutes.HandleFunc("/{id}", company.DeleteCompany).Methods("DELETE")

	projectsRoutes := r.PathPrefix("/projects").Subrouter()
	projectsRoutes.Use(core.AuthMiddleware)
	projectsRoutes.Use(core.AuthorizationMiddleware)
	projectsRoutes.HandleFunc("", projects.GetProjects).Methods("GET")
	projectsRoutes.HandleFunc("/{id}", projects.GetProject).Methods("GET")
	projectsRoutes.HandleFunc("", projects.CreateProject).Methods("POST")
	projectsRoutes.HandleFunc("/{id}", projects.UpdateProject).Methods("PUT")
	projectsRoutes.HandleFunc("/{id}", projects.DeleteProject).Methods("DELETE")

	recordsRoutes := r.PathPrefix("/records").Subrouter()
	recordsRoutes.Use(core.AuthMiddleware)
	recordsRoutes.Use(core.AuthorizationMiddleware)
	recordsRoutes.HandleFunc("", records.GetRecords).Methods("GET")
	recordsRoutes.HandleFunc("/{id}", records.GetRecord).Methods("GET")
	recordsRoutes.HandleFunc("", records.CreateRecord).Methods("POST")
	recordsRoutes.HandleFunc("/{id}", records.UpdateRecord).Methods("PUT")
	recordsRoutes.HandleFunc("/{id}", records.DeleteRecord).Methods("DELETE")

	r.HandleFunc("/auth/login", users.Login).Methods("POST")
	r.HandleFunc("/auth/sign-up", users.CreateUser).Methods("POST")

	return r
}

func InitEnv() {
	env := os.Getenv("ENV")
	if env == "" {
		env = "local"
	}

	godotenv.Load(".env." + env)
	err := godotenv.Load() // The Original .env
	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}
}

func main() {
	InitEnv()
	core.InitDB()
	core.InitValidation()
	core.InitAuthorization()
	r := initRouter()

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	log.Printf("Listening on %s:%s\n", host, port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)

	if err != nil {
		log.Fatal(err.Error())
	}
}
