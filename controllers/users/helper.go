package users

import "github.com/mfjkri/OneNUS-Backend/models"

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
