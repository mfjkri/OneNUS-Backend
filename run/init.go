package run

import "github.com/joho/godotenv"

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}
}