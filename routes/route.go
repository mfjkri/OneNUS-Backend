package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/auth/register", controllers.RegisterJSON)
	r.POST("/auth/login", controllers.LoginJSON)
	r.GET("/auth/me", controllers.GetUserJSON)
	
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
}