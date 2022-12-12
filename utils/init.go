package utils

import (
	"fmt"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}

	fmt.Println("Successfully import .env variables...")
}

