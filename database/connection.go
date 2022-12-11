package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	connection, err := gorm.Open(postgres.Open(os.Getenv("DB")) , &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	DB = connection
}