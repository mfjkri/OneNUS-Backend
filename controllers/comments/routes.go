package comments

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.GET("comments/get/:postId/:perPage/:pageNumber/:sortOption/:sortOrder", GetComments)
	r.POST("comments/create", CreateComment)
	r.POST("comments/updatetext", UpdateCommentText)
	r.DELETE("comments/delete/:commentId", DeleteComment)
}
