package database

import (
	"fmt"

	"github.com/mfjkri/OneNUS-Backend/models"
)

func Migrate() {
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

	fmt.Println("Sucessfully migrated database...")
}
