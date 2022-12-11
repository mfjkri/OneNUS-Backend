package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	fmt.Println(os.Getenv("DB"))
	connection, err := gorm.Open(postgres.Open(os.Getenv("DB")) , &gorm.Config{})

	if err != nil {
		panic("Failed to connect to database!")
	}

	fmt.Println("Successfully connected to database...")

	DB = connection
}