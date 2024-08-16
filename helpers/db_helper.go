package helpers

import (
	"e_commerce/database"

	"gorm.io/gorm"
)

func GetDB() *gorm.DB {
	return database.DB
}
