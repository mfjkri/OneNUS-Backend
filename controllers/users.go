package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/models"
)

type UserResponse struct {
	ID            uint   `json:"id" binding:"required"`
	Username      string `json:"username" binding:"required"`
	Role          string `json:"role" binding:"required"`
	Bio           string `json:"bio" binding:"required"`
	PostsCount    uint   `json:"postsCount" binding:"required"`
	CommentsCount uint   `json:"commentsCount" binding:"required"`
	CreatedAt     int64  `json:"createdAt" binding:"required"`
}

func CreateUserResponse(user *models.User) UserResponse {
	return UserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Role:          user.Role,
		Bio:           user.Bio,
		PostsCount:    user.PostsCount,
		CommentsCount: user.CommentsCount,
		CreatedAt:     user.CreatedAt.Unix(),
	}
}

/* -------------------------------------------------------------------------- */
/*                GetUserFromID | route: /users/getbyid/:userId               */
/* -------------------------------------------------------------------------- */
type GetUserFromIDRequest struct {
	UserID uint `uri:"userId" binding:"required"`
}

func GetUserFromID(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetUserFromIDRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user, found := FindUserFromID(c, json.UserID)

	if found == false {
		return
	}

	// Return fetch User
	c.JSON(http.StatusAccepted, CreateUserResponse(&user))
}
