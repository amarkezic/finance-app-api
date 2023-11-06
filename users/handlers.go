package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"amarkezic.github.com/finance-app/core"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

var GetUsers = core.List[core.User]()

var GetUser = core.Single[core.User]()

var CreateUser = core.Create[core.User]()

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
	token, err := core.GenerateToken()

	if err != nil {
		response := core.AuthResponse{
			Ok:      false,
			Message: err.Error(),
			Data:    nil,
		}
		json.NewEncoder(w).Encode(response)
		return
	}

	response := core.AuthResponse{
		Ok:      true,
		Message: "",
		Data:    token,
	}

	json.NewEncoder(w).Encode(response)
}
