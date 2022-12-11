package controllers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
	"golang.org/x/crypto/bcrypt"
)

type Register struct {
	Username	string 	`form:"username" json:"username" binding:"required"`
	Password	string 	`form:"password" json:"password" binding:"required"`
}

func RegisterJSON(c *gin.Context) {
	// Parse RequestBody 
	var json Register
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
	jwt, err := GenerateJWT(username_lowered)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": "Failed to create access token."})
      	return
	}

	// Success, registered and logged in
    CreateUserResponse(c, http.StatusOK, jwt, user.ID, user.Username)
}

type Login struct {
	Username	string 	`form:"username" json:"username" binding:"required"`
	Password 	string 	`form:"password" json:"password" binding:"required"`
}

func LoginJSON(c *gin.Context) {
	// Parse RequestBody 	
	var json Login
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

	jwt, err := GenerateJWT(username_lowered)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, gin.H{"message": "Failed to create access token."})
      	return
	}

	// Success, logged in
    CreateUserResponse(c, http.StatusOK, jwt, user.ID, user.Username)
}

func GetUserJSON(c *gin.Context) {
	user, found := RequireAuth(c)

	if found == false {
		return
	}

	// Success, user found
	CreateUserResponseNoJWT(c, http.StatusAccepted, user.Id, user.Username)
}