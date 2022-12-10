package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mfjkri/One-NUS-Backend/controllers"
)

func Setup(app *fiber.App) {
	app.Post("/auth/register", controllers.Register)
	app.Post("/auth/login", controllers.Login)
	app.Get("/auth/me", controllers.User)
}