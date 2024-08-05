package database

import (
	"e_commerce/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := os.Getenv("DSL")
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	DB = db
}

func Migrate() {
	DB.AutoMigrate(&models.User{})
	DB.AutoMigrate(&models.Shop{})
	DB.AutoMigrate(&models.Product{})
}
