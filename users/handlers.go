package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"amarkezic.github.com/finance-app/core"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []core.User
	core.DB.Find(&users)
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user core.User
	result := core.DB.First(&user, params["id"])

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user core.User
	json.NewDecoder(r.Body).Decode(&user)

	validationResult := core.ValidateStruct(user)

	if len(validationResult) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(validationResult)
		return
	}

	hashedPassword, err := core.HashPassword(user.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user.Password = hashedPassword
	core.DB.Create(&user)
	json.NewEncoder(w).Encode(user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user core.User
	core.DB.First(&user, params["id"])
	json.NewDecoder(r.Body).Decode(&user)

	validationResult := core.ValidateStruct(user)
	if len(validationResult) > 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(validationResult)
		return
	}

	core.DB.Save(&user)
	json.NewEncoder(w).Encode(user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user core.User
	core.DB.Delete(&user, params["id"])
	w.WriteHeader(http.StatusOK)
}
