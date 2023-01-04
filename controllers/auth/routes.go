package auth

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.POST("auth/register", RegisterUser)
	r.POST("auth/login", LoginUser)
	r.GET("auth/me", GetUser)
}
