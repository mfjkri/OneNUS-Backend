package posts

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.GET("posts/get/:perPage/:pageNumber/:sortOption/:sortOrder/:filterUserId/:filterTag", GetPosts)
	r.GET("posts/getbyid/:postId", GetPostByID)
	r.POST("posts/create", CreatePost)
	r.POST("posts/updatetext", UpdatePostText)
	r.DELETE("posts/delete/:postId", DeletePost)
}
