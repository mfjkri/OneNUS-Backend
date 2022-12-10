package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/routes"
)


func main() {
	database.Connect()

	app := fiber.New()

	routes.Setup(app)

	app.Listen(":8000")
}