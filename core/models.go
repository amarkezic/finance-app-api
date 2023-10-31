package core

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8"`
	DisplayName string `json:"displayName" validate:"required"`
}
