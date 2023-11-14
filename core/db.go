package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func InitDB() {
	connectionString := os.Getenv("DB_CONNECTION_STRING")
	DB, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot establish a connection to the database")
	}

	DB.AutoMigrate(&User{})
	DB.AutoMigrate(&Record{})
	DB.AutoMigrate(&Project{})
	DB.AutoMigrate(&Company{})
}

func List[T any]() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var entities []T
		DB.Find(&entities)
		json.NewEncoder(w).Encode(entities)
	}
}

func Single[T any]() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)

		fieldsToPopulate := strings.Split(r.URL.Query().Get("populate"), ",")

		dbModel := DB

		for _, value := range fieldsToPopulate {
			dbModel = dbModel.Preload(cases.Title(language.English).String(value))
		}

		var entity T
		result := dbModel.First(&entity, params["id"])

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewEncoder(w).Encode(entity)
	}
}

func Create[T any]() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var entity T
		json.NewDecoder(r.Body).Decode(&entity)

		validationResult := ValidateStruct(entity)

		if len(validationResult) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(validationResult)
			return
		}

		DB.Create(&entity)
		json.NewEncoder(w).Encode(entity)
	}
}

func Update[T any]() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		params := mux.Vars(r)
		var entity T
		result := DB.First(&entity, params["id"])

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		json.NewDecoder(r.Body).Decode(&entity)

		validationResult := ValidateStruct(entity)
		if len(validationResult) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(validationResult)
			return
		}

		DB.Save(&entity)
		json.NewEncoder(w).Encode(entity)
	}
}

func Delete[T any]() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)
		var entity T

		result := DB.Delete(&entity, params["id"])

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
