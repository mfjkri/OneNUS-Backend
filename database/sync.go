package database

import (
	"fmt"

	"github.com/mfjkri/One-NUS-Backend/models"
)

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	fmt.Println("Sucessfully migrated database...")
}