package core

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const connectionString = "root:Admin.1234@tcp(127.0.0.1:3306)/finance-app?charset=utf8mb4&parseTime=True&loc=Local"

func InitDB() {
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
