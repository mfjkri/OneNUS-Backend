package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/controllers"
)

func RegisterPublicRoutes(r *gin.Engine) {
	// auth.go
	r.POST("auth/register", controllers.RegisterUser)
	r.POST("auth/login", controllers.LoginUser)

	// misc
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "oneNUS " + os.Getenv("APP_VERSION"),
		})
	})
}
