package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/register", controllers.RegisterJSON)
	r.POST("/auth/login", controllers.LoginJSON)
}