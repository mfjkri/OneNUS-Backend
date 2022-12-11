package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
)

func GenerateJWT(username string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub" : username,
		"exp" : time.Now().Add(time.Hour * 8760).Unix(),
	})

	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func DecodeJWT(tokenString string) (username string, err error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username = claims["sub"].(string)
		return 
	}

	username = "" 
	return 
}

func ValidateJWT(tokenString string, username string) (jwtValid bool, err error) {
	decoded_username, err := DecodeJWT(tokenString)

	if decoded_username == username {
		jwtValid = true
		return 
	}

	jwtValid = false
	return 
}

type UserAuth struct {
	Id     		uint 	`json:"id" binding:"required"`
	Username 	string 	`json:"username" binding:"required"`
}

type UserResponse struct {
	JWT     string 		`json:"jwt" binding:"required"`
	User 	UserAuth 	`json:"user" binding:"required"`
}

func CreateUserResponse(c *gin.Context, http_status uint, jwt string, id uint, username string) {
	c.JSON(int(http_status), gin.H{
		"status": "Success",
		"jwt": jwt,
		"user": gin.H{
			"id": id,
			"username": username,
		},
	})
}

func CreateUserResponseNoJWT(c *gin.Context, http_status uint, id uint, username string) {
	c.JSON(int(http_status), gin.H{
		"id": id,
		"username": username,
	})
}


func SanitizeUser(user UserAuth) UserAuth{
	return UserAuth{
		user.Id,
		user.Username,
	}
}

func RequireAuth(c *gin.Context) (user UserAuth, found bool) {
	jwt_token := c.Request.Header.Get("authorization");
	found = false

	if jwt_token == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "No authorization token provided."})
      	return 
	}

	username, err := DecodeJWT(jwt_token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to decoded authorization token."})
      	return 
	}

	var target_user models.User
    database.DB.First(&target_user, "username = ?", strings.ToLower(username))

	if target_user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unauthorized."})
      	return 
	}

	found = true
	user = UserAuth{
		Id : target_user.ID,
		Username:  target_user.Username,
	}
	return

}