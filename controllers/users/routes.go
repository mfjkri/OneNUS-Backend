package users

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine) {
	r.GET("users/getbyid/:userId", GetUserFromID)
	r.POST("users/updatebio", UpdateBio)
	r.DELETE("users/delete", DeleteUser)
}
