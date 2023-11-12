package users

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"amarkezic.github.com/finance-app/core"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var GetUsers = core.List[core.User]()

var GetUser = core.Single[core.User]()

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var userBody core.User
	json.NewDecoder(r.Body).Decode(&userBody)

	validationResult := core.ValidateStruct(userBody)

	if len(validationResult) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationResult)
		return
	}

	hashedPassword, err := core.HashPassword(userBody.Password)

	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userBody.Password = hashedPassword

	core.DB.Create(&userBody)
	json.NewEncoder(w).Encode(userBody)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var user core.User
	result := core.DB.First(&user, params["id"])

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

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

var DeleteUser = core.Delete[core.User]()

func Login(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var userBody core.User

	json.NewDecoder(r.Body).Decode(&userBody)

	var user core.User
	userQueryResult := core.DB.First(&user, "email = ?", strings.ToLower(userBody.Email))

	if errors.Is(userQueryResult.Error, gorm.ErrRecordNotFound) {
		response := core.Response{
			Ok:      false,
			Message: "Wrong username of password",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	passwordValid := core.CheckPasswordValidity(user.Password, userBody.Password)

	if !passwordValid {
		response := core.Response{
			Ok:      false,
			Message: "Wrong username of password",
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	token, err := core.GenerateToken(&user)

	if err != nil {
		response := core.AuthResponse{
			Ok:      false,
			Message: err.Error(),
			Token:   nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := core.AuthResponse{
		Ok:      true,
		Message: "",
		Token:   token,
	}

	json.NewEncoder(w).Encode(response)
}
