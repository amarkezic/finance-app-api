package company

import (
	"encoding/json"
	"errors"
	"net/http"

	"amarkezic.github.com/finance-app/core"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func GetCompanies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var companies []core.Company
	core.DB.Find(&companies)
	json.NewEncoder(w).Encode(companies)
}

func GetCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var company core.Company
	result := core.DB.Preload("Users").Preload("Projects").First(&company, params["id"])

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(company)
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var company core.Company
	json.NewDecoder(r.Body).Decode(&company)

	validationResult := core.ValidateStruct(company)

	if len(validationResult) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationResult)
		return
	}

	core.DB.Create(&company)
	json.NewEncoder(w).Encode(company)
}

func UpdateCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var company core.Company
	result := core.DB.First(&company, params["id"])

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	json.NewDecoder(r.Body).Decode(&company)

	validationResult := core.ValidateStruct(company)
	if len(validationResult) > 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(validationResult)
		return
	}

	core.DB.Save(&company)
	json.NewEncoder(w).Encode(company)
}

func DeleteCompany(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var company core.Company

	result := core.DB.Delete(company, params["id"])

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
}
