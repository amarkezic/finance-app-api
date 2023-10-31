package users

import (
	"encoding/json"
	"net/http"

	"amarkezic.github.com/finance-app/core"
)

func GetUsers(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Content-Type", "application/json")
	var users []core.User
	core.DB.Find(&users)
	json.NewEncoder(response).Encode(users)
}

func GetUser(response http.ResponseWriter, request *http.Request) {

}

func CreateUser(response http.ResponseWriter, request *http.Request) {

}

func UpdateUser(response http.ResponseWriter, request *http.Request) {

}
