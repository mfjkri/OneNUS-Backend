package database

import "github.com/mfjkri/One-NUS-Backend/models"

func Migrate() {
	DB.AutoMigrate(&models.User{})
}