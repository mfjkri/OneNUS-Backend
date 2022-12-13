package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/controllers"
)

func SetupRoutes(r *gin.Engine) {
	// auth.go
	r.POST("auth/register", controllers.RegisterUser)
	r.POST("auth/login", controllers.LoginUser)
	r.GET("auth/me", controllers.GetUser)

	// posts.go
	r.POST("posts/create", controllers.CreatePost)
	r.GET("posts/get/:perPage/:pageNumber/:sortBy/:filterTag", controllers.GetPosts)
	r.GET("posts/getbyid", controllers.GetPostsByID)
	r.POST("/posts/updatetext", controllers.UpdatePostText)
	
	// misc
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 1",
		})
	})
}