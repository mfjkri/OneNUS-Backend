package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"github.com/mfjkri/OneNUS-Backend/utils"
)

func FindUserFromID(c *gin.Context, userID uint) (models.User, bool) {
	var targetUser models.User
	database.DB.First(&targetUser, userID)
	if targetUser.ID == 0 {
		c.JSON(http.StatusForbidden, gin.H{"message": "User not found."})
		return targetUser, false
	} else {
		return targetUser, true
	}
}

// Verify RequestUser using their JWT token
func VerifyAuth(c *gin.Context) (user models.User, found bool) {
	found = false
	jwt_token := c.Request.Header.Get("authorization")

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

// Convert a User Model into a JSON format
type AuthResponse struct {
	ID       uint   `json:"id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

func CreateAuthResponse(user *models.User) AuthResponse {
	return AuthResponse{
		ID:       user.ID,
		Username: user.Username,
		Role:     user.Role,
	}
}

// Convert a User Model with JWT into a JSON format
type AuthResponseWithJWT struct {
	JWT  string       `json:"jwt" binding:"required"`
	User AuthResponse `json:"user" binding:"required"`
}

func CreateAuthResponseWithJWT(jwt string, user *models.User) AuthResponseWithJWT {
	return AuthResponseWithJWT{
		JWT:  jwt,
		User: CreateAuthResponse(user),
	}
}
