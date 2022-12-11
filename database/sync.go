package database

import (
	"fmt"

	"github.com/mfjkri/One-NUS-Backend/models"
)

func Migrate() {
	DB.AutoMigrate(&models.User{})

	fmt.Println("Sucessfully migrated database...")
}