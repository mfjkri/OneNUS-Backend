package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
	"github.com/mfjkri/One-NUS-Backend/utils"
	"golang.org/x/crypto/bcrypt"
)



type UserAuth struct {
	Id     		uint 	`json:"id" binding:"required"`
	Username 	string 	`json:"username" binding:"required"`
}


func VerifyAuth(c *gin.Context) (user UserAuth, found bool) {
	jwt_token := c.Request.Header.Get("authorization");
	found = false

	if jwt_token == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "No authorization token provided."})
      	return 
	}

	
	username, err := utils.DecodeJWT(jwt_token)
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

type RegisterRequest struct {
	Username	string 	`form:"username" json:"username" binding:"required"`
	Password	string 	`form:"password" json:"password" binding:"required"`
}

func RegisterUser(c *gin.Context) {
	// Parse RequestBody 
	var json RegisterRequest
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
      return
    }

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(json.Password), 10)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "Failed to hash password."})
      	return
	}

	username_lowered := strings.ToLower(json.Username)

    user := models.User{Username: username_lowered, Password: hash}
	new_entry := database.DB.Create(&user)

	// Failed to create entry: most likely user already exists
	if new_entry.Error != nil {
		c.JSON(http.StatusConflict, gin.H{"message": "User already exists."})
      	return
	}

	// Generate JWT Token
	jwt, err := utils.GenerateJWT(username_lowered)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": "Failed to create access token."})
      	return
	}

	// Success, registered and logged in
    CreateUserResponse(c, http.StatusOK, jwt, user.ID, user.Username)
}

type LoginRequest struct {
	Username	string 	`form:"username" json:"username" binding:"required"`
	Password 	string 	`form:"password" json:"password" binding:"required"`
}

func LoginUser(c *gin.Context) {
	// Parse RequestBody 	
	var json LoginRequest
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
      return
    }

	username_lowered := strings.ToLower(json.Username)

	// Find User based on request.username
	var user models.User
    database.DB.First(&user, "username = ?", username_lowered)
	if user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Username not found."})
      	return
	}

	// Compare password and saved hash
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(json.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Wrong password. Please try again."})
      	return
	}

	jwt, err := utils.GenerateJWT(username_lowered)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": "Failed to create access token."})
      	return
	}

	// Success, logged in
    CreateUserResponse(c, http.StatusOK, jwt, user.ID, user.Username)
}

func GetUser(c *gin.Context) {
	user, found := VerifyAuth(c)

	if found == false {
		return
	}

	// Success, user found
	c.JSON(http.StatusAccepted, gin.H{
		"id": user.Id,
		"username": user.Username,
	})
}