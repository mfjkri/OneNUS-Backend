package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("auth/register", controllers.RegisterUser)
	r.POST("auth/login", controllers.LoginUser)
	r.GET("auth/me", controllers.GetUser)
	
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 1",
		})
	})
}