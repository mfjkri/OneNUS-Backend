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
	// r.POST("auth/delete", controllers.DeleteUser)

	// posts.go
	r.GET("posts/get/:perPage/:pageNumber/:sortBy/:filterTag", controllers.GetPosts)
	r.GET("posts/getbyid/:postId", controllers.GetPostByID)
	r.POST("posts/create", controllers.CreatePost)
	r.POST("posts/updatetext", controllers.UpdatePostText)
	r.DELETE("posts/delete/:postId", controllers.DeletePost)

	// comments.go
	r.GET("comments/get/:postId/:perPage/:pageNumber/:sortBy", controllers.GetComments)
	r.POST("comments/create", controllers.CreateComment)
	r.POST("comments/updatetext", controllers.UpdateCommentText)
	r.DELETE("comments/delete/:commentId", controllers.DeleteComment)
	
	// misc
	r.GET("ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong 1",
		})
	})
}