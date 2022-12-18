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

/* -------------------------------------------------------------------------- */
/*                              Helper functions                              */
/* -------------------------------------------------------------------------- */

// Verify RequestUser using their JWT token
func VerifyAuth(c *gin.Context) (user models.User, found bool) {
	found = false
	jwt_token := c.Request.Header.Get("authorization");

	// Check for authorization token (JWT)
	if jwt_token == "" {
		c.JSON(http.StatusNotFound, gin.H{"message": "No authorization token provided."})
      	return 
	}
	
	// Decode JWT token to username
	username, err := utils.DecodeJWT(jwt_token)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unable to decoded authorization token."})
      	return 
	}

	// Search for User from username
	var target_user models.User
    database.DB.Table("users").Where("username = ?", strings.ToLower(username)).First(&target_user)
	if target_user.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Unauthorized."})
      	return 
	}

	// User successfully found
	found = true
	user = target_user
	return
}

// Generate JSONResponse with jwt, int, username
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
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                    RegisterUser | route: /auth/register                    */
/* -------------------------------------------------------------------------- */
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

	// Check that Username and Password does not contain illegal characters
	if utils.ContainsWhitespaces(json.Password) || !utils.ContainsLettersOnly(json.Username) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Username or password contains illegal characters."})
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
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                       LoginUser | route: /auth/login                       */
/* -------------------------------------------------------------------------- */
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
	database.DB.Table("users").Where("username = ?", username_lowered).First(&user)
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
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                          GetUser | route: /auth/me                         */
/* -------------------------------------------------------------------------- */
func GetUser(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Success, user found
	c.JSON(http.StatusAccepted, gin.H{
		"id": user.ID,
		"username": user.Username,
	})
}
/* -------------------------------------------------------------------------- */


/* -------------------------------------------------------------------------- */
/*                      DeleteUser | route : /auth/delete                     */
/* -------------------------------------------------------------------------- */
func DeleteUser(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	database.DB.Delete(&user)

	// Success, user deleted
	c.JSON(http.StatusAccepted, gin.H{
		"username": user.Username,
	})
}
/* -------------------------------------------------------------------------- */