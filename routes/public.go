package routes

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/controllers/auth"
)

func RegisterPublicRoutes(r *gin.Engine) {
	// auth.go
	auth.RegisterRoutes(r)

	// misc
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "oneNUS " + os.Getenv("APP_VERSION"),
		})
	})
}
