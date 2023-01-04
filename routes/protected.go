package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/controllers/comments"
	"github.com/mfjkri/OneNUS-Backend/controllers/posts"
	"github.com/mfjkri/OneNUS-Backend/controllers/users"
)

func RegisterProtectedRoutes(r *gin.Engine) {
	posts.RegisterRoutes(r)
	comments.RegisterRoutes(r)
	users.RegisterRoutes(r)
}
