package controllers

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func GenerateJWT(id uint) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub" : id,
		"exp" : time.Now().Add(time.Hour * 8760).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

type UserAuth struct {
	Id     		string 	`json:"id" binding:"required"`
	Username 	string 	`json:"username" binding:"required"`
}

type UserResponse struct {
	JWT     string 		`json:"jwt" binding:"required"`
	User 	UserAuth 	`json:"user" binding:"required"`
}

func CreateUserResponse(c *gin.Context, http_status uint, jwt string, id uint, username string) {
	c.JSON(int(http_status), gin.H{
		"status": "you are logged in",
		"jwt": jwt,
		"user": gin.H{
			"id": id,
			"username": username,
		},
	})
}