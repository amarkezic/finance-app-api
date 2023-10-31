package core

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email       string `json:"email"`
	Password    string `json:"password"`
	DisplayName string `json:"DisplayName"`
}