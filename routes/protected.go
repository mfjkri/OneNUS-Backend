package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/controllers"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	// auth.go
	r.GET("auth/me", controllers.GetUser)

	// posts.go
	r.GET("posts/get/:perPage/:pageNumber/:sortOption/:sortOrder/:filterUserId/:filterTag", controllers.GetPosts)
	r.GET("posts/getbyid/:postId", controllers.GetPostByID)
	r.POST("posts/create", controllers.CreatePost)
	r.POST("posts/updatetext", controllers.UpdatePostText)
	r.DELETE("posts/delete/:postId", controllers.DeletePost)

	// comments.go
	r.GET("comments/get/:postId/:perPage/:pageNumber/:sortOption/:sortOrder", controllers.GetComments)
	r.POST("comments/create", controllers.CreateComment)
	r.POST("comments/updatetext", controllers.UpdateCommentText)
	r.DELETE("comments/delete/:commentId", controllers.DeleteComment)

	// users.go
	r.GET("users/getbyid/:userId", controllers.GetUserFromID)
	r.POST("users/updatebio", controllers.UpdateBio)
	r.POST("users/delete", controllers.DeleteUser)
}
