package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
	"github.com/mfjkri/One-NUS-Backend/utils"
)

/* -------------------------------------------------------------------------- */
/*                              Helper functions                              */
/* -------------------------------------------------------------------------- */
type CommentResponse struct {
	ID     uint   `json:"id" binding:"required"`
	Text   string `json:"text" binding:"required"`
	Author string `json:"author" binding:"required"`
	UserID uint   `json:"userId" binding:"required"`

	PostID uint `json:"postId" binding:"required"`

	CreatedAt int64 `json:"createdAt" binding:"required"`
	UpdatedAt int64 `json:"updatedAt" binding:"required"`
}

// Convert a Comment Model into a JSON format
func CreateCommentResponse(comment *models.Comment) CommentResponse {
	return CommentResponse{
		ID:        comment.ID,
		Text:      comment.Text,
		Author:    comment.Author,
		UserID:    comment.UserID,
		PostID:    comment.PostId,
		CreatedAt: comment.CreatedAt.Unix(),
		UpdatedAt: comment.UpdatedAt.Unix(),
	}
}

type GetCommentsResponse struct {
	Comments      []CommentResponse `json:"comments" binding:"required"`
	CommentsCount int64             `json:"commentsCount" binding:"required"`
}

// Bundles and convert multiple comments models into a JSON format
func CreateCommentsResponse(comments *[]models.Comment, totalCommentsCount int64) GetCommentsResponse {
	var commentsResponse []CommentResponse
	for _, comment := range *comments {
		commentResponse := CreateCommentResponse(&comment)
		commentsResponse = append(commentsResponse, commentResponse)
	}

	return GetCommentsResponse{
		Comments:      commentsResponse,
		CommentsCount: totalCommentsCount,
	}
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*   GetComments | route: comments/get/:postId/:perPage/:pageNumber/:sortBy   */
/* -------------------------------------------------------------------------- */
type GetCommentsRequest struct {
	PostId     uint   `uri:"postId" binding:"required"`
	PerPage    uint   `uri:"perPage" binding:"required"`
	PageNumber uint   `uri:"pageNumber" binding:"required"`
	SortBy     string `uri:"sortBy"`
}

func GetComments(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetCommentsRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Find Post from PostId
	var post models.Post
	database.DB.First(&post, json.PostId)
	if post.ID == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "Post not found."})
		return
	}

	// Limit PerPage to MAX_PER_PAGE
	perPage := int64(math.Min(MAX_PER_PAGE, float64(json.PerPage)))
	offsetCommentCount := int64(json.PageNumber-1) * perPage

	// Get all comments from Post
	dbContext := database.DB.Table("comments").Where("post_id = ?", json.PostId)

	// Get total count for Comments
	totalCommentsCount := int64(post.CommentsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetCommentCount < 0) || (offsetCommentCount > totalCommentsCount) {
		c.JSON(http.StatusForbidden, gin.H{"message": "No more comments found."})
		return
	}

	// Sort Comments by sort option provided (defaults to byNew)
	defaultSortOption := ByNew
	if json.SortBy == "byRecent" {
		defaultSortOption = ByRecent
	}

	// Fetch Posts from [offsetCount, offsetCount + perPage]
	var comments []models.Comment
	dbContext.Limit(int(perPage)).Order(defaultSortOption).Offset(int(offsetCommentCount)).Find(&comments)

	// Return fetched posts
	c.JSON(http.StatusAccepted, CreateCommentsResponse(&comments, totalCommentsCount))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                   CreateComment | route: comments/create                   */
/* -------------------------------------------------------------------------- */
type CreateCommentRequest struct {
	PostId uint   `json:"postId" binding:"required"`
	Text   string `json:"text" binding:"required"`
}

func CreateComment(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json CreateCommentRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Find Post from PostId
	var post models.Post
	database.DB.First(&post, json.PostId)
	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	// Prevent frequent CreatePosts by User
	timeNow, canCreateComment := utils.CheckTimeIsAfter(user.LastCommentAt, USER_COMMENT_COOLDOWN)
	if canCreateComment == false {
		cdLeft := utils.GetCooldownLeft(user.LastCommentAt, USER_COMMENT_COOLDOWN, timeNow)
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Creating comments too frequently. Please try again in %ds", int(cdLeft.Seconds()))})
		return
	}

	// Check that Title and Text does not contain illegal characters
	if !(utils.ContainsValidCharactersOnly(json.Text)) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Comment contains illegal characters."})
		return
	}

	// Try to create new Comment
	comment := models.Comment{
		Text:   utils.TrimString(strings.TrimSpace(json.Text), MAX_COMMENT_TEXT_CHAR),
		Author: user.Username,
		User:   user,
		Post:   post,
	}
	new_entry := database.DB.Create(&comment)

	// Failed to create entry
	if new_entry.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Unable to create comment. Try again later."})
		return
	}

	// Successfully created a new Post

	// Update LastCommentAt for User
	user.LastCommentAt = timeNow
	database.DB.Save(&user)

	// Update CommentsCount for Post
	post.CommentedAt = timeNow
	post.CommentsCount += 1
	database.DB.Save(&post)

	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*               UpdateCommentText | route: comments/updatetext               */
/* -------------------------------------------------------------------------- */
type UpdateCommentTextRequest struct {
	Text      string `json:"text" binding:"required"`
	CommentId uint   `json:"commentId" binding:"required"`
}

func UpdateCommentText(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json UpdateCommentTextRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Find Comment from CommentId
	var comment models.Comment
	database.DB.First(&comment, json.CommentId)
	if comment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	}

	// Prevent frequent UpdateCommentText by User
	timeNow, canUpdateComment := utils.CheckTimeIsAfter(user.LastCommentAt, USER_COMMENT_COOLDOWN)
	if canUpdateComment == false {
		cdLeft := utils.GetCooldownLeft(user.LastCommentAt, USER_COMMENT_COOLDOWN, timeNow)
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Updating comments too frequently. Please try again in %ds", int(cdLeft.Seconds()))})
		return
	}

	// Check User is the author
	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Replace Comment text and update User LastCommentAt
	comment.Text = utils.TrimString(strings.TrimSpace(json.Text), MAX_COMMENT_TEXT_CHAR)
	user.LastCommentAt = timeNow
	database.DB.Save(&comment)
	database.DB.Save(&user)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*              DeleteComment | route: comments/delete/:commentId             */
/* -------------------------------------------------------------------------- */
type DeleteCommentRequest struct {
	CommentId uint `uri:"commentId" binding:"required"`
}

func DeleteComment(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json DeleteCommentRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Find Comment from CommentId
	var comment models.Comment
	database.DB.First(&comment, json.CommentId)
	if comment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	}

	// Find Post from PostId
	var post models.Post
	database.DB.First(&post, comment.PostId)
	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	// Check User is the author or is admin
	if (comment.UserID != user.ID) && (user.Role != ADMIN) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Delete Comment
	database.DB.Delete(&comment)

	// Update Post metadata
	var lastComment models.Comment
	database.DB.Table("comments").Where("post_id = ?", post.ID).Order("created_at DESC, id DESC").First(&lastComment)
	if lastComment.ID != 0 {
		post.CommentedAt = lastComment.CreatedAt
	} else {
		post.CommentedAt = time.Unix(0, 0)
	}
	post.CommentsCount -= 1
	database.DB.Save(&post)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}

/* -------------------------------------------------------------------------- */
