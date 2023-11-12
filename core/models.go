package core

import (
	"gorm.io/gorm"
)

type ProjectStatus int

const (
	Created   ProjectStatus = 0
	Active    ProjectStatus = 1
	Finished  ProjectStatus = 2
	Completed ProjectStatus = 3
)

type RecordStatus int

const (
	Pending RecordStatus = 0
	Paid    RecordStatus = 1
)

type UserRole string

const (
	Admin   UserRole = "Admin"
	Default UserRole = "Default"
)

type User struct {
	gorm.Model
	Email       string   `json:"email" validate:"required,email"`
	Password    string   `json:"password" validate:"required,min=8"`
	DisplayName string   `json:"displayName" validate:"required"`
	Role        UserRole `json:"role" validate:"required"`
	CompanyId   uint     `json:"companyId"`
}

type Project struct {
	gorm.Model
	CompanyId uint          `json:"companyId" validate:"required"`
	Name      string        `json:"name" validate:"required"`
	Status    ProjectStatus `json:"status" validate:"required"`
	Records   []Record      `json:"records"`
}

type Company struct {
	gorm.Model
	Name     string    `json:"name" validate:"required"`
	Projects []Project `json:"projects"`
	Users    []User    `json:"users"`
}

type Record struct {
	gorm.Model
	ProjectId uint         `json:"projectId" validate:"required"`
	Status    RecordStatus `json:"status" validate:"required"`
	Value     float64      `json:"value" validate:"required"`
	Tax       float64      `json:"tax" validate:"required"`
}

type Response struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}

type AuthResponse struct {
	Ok      bool        `json:"ok"`
	Message string      `json:"message"`
	Token   interface{} `json:"token"`
}
