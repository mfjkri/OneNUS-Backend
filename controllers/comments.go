package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"github.com/mfjkri/OneNUS-Backend/utils"
	"gorm.io/gorm"
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
		PostID:    comment.PostID,
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

// Fetches comments based on provided configuration
func GetCommentsFromContext(dbContext *gorm.DB, perPage uint, pageNumber uint, sortOption string, sortOrder string) ([]models.Comment, int64) {
	var comments []models.Comment

	// Limit PerPage to MAX_PER_PAGE
	clampedPerPage := int64(math.Min(MAX_PER_PAGE, float64(perPage)))
	offsetCommentsCount := int64(pageNumber-1) * clampedPerPage

	// Get total count for Comments
	var totalCommentsCount int64
	dbContext.Count(&totalCommentsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetCommentsCount < 0) || (offsetCommentsCount > totalCommentsCount) {
		return comments, 0
	}

	// Sort Comments by sort option provided (defaults to byNew)
	defaultSortOption := ByNew
	if sortOption == "recent" {
		defaultSortOption = ByRecent
	}

	// Fetch Comments from [offsetCount, offsetCount + perPage]
	// results order depends on SortOption and SortOrder
	if sortOrder == "ascending" {
		// Reverse page number based on totalPostsCount
		leftOverRecords := math.Min(float64(perPage), float64(totalCommentsCount-offsetCommentsCount))
		offsetCommentsCount = totalCommentsCount - offsetCommentsCount - clampedPerPage
		dbContext.Limit(int(leftOverRecords)).Order(defaultSortOption).Offset(int(offsetCommentsCount)).Find(&comments)

		// Reverse the page results for descending order
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}
	} else {
		dbContext.Limit(int(perPage)).Order(defaultSortOption).Offset(int(offsetCommentsCount)).Find(&comments)
	}

	return comments, totalCommentsCount
}

/* -------------------------------------------------------------------------- */
/*   GetComments | route: comments/get/:postId/:perPage/:pageNumber/:sortBy   */
/* -------------------------------------------------------------------------- */
type GetCommentsRequest struct {
	PostID     uint   `uri:"postId" binding:"required"`
	PerPage    uint   `uri:"perPage" binding:"required"`
	PageNumber uint   `uri:"pageNumber" binding:"required"`
	SortOption string `uri:"sortOption"`
	SortOrder  string `uri:"sortOrder"`
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

	// Find Post from PostID
	var post models.Post
	database.DB.First(&post, json.PostID)
	if post.ID == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "Post not found."})
		return
	}

	// Get all comments from Post
	dbContext := database.DB.Table("comments").Where("post_id = ?", json.PostID)
	comments, totalCommentsCount := GetCommentsFromContext(dbContext, json.PerPage, json.PageNumber, json.SortOption, json.SortOrder)

	// Return fetched comments
	c.JSON(http.StatusAccepted, CreateCommentsResponse(&comments, totalCommentsCount))
}

/* -------------------------------------------------------------------------- */
/*                   CreateComment | route: comments/create                   */
/* -------------------------------------------------------------------------- */
type CreateCommentRequest struct {
	PostID uint   `json:"postId" binding:"required"`
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

	// Find Post from PostID
	var post models.Post
	database.DB.First(&post, json.PostID)
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

	// Update CommentsCount and CommentedAt for Post
	postUpdatedAt := post.UpdatedAt
	post.CommentedAt = timeNow
	post.CommentsCount += 1
	database.DB.Save(&post)

	// By default post.UpdatedAt will get updated by these changes but we don't want that.
	// UpdatedAt should only reflect changes to post.Text
	database.DB.Model(&post).Update("updated_at", postUpdatedAt)

	fmt.Printf("%s has created a comment.\n\tPost title: %s\n\tComment text: %s\n", user.Username, post.Title, comment.Text)

	// Return new Comment data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}

/* -------------------------------------------------------------------------- */
/*               UpdateCommentText | route: comments/updatetext               */
/* -------------------------------------------------------------------------- */
type UpdateCommentTextRequest struct {
	Text      string `json:"text" binding:"required"`
	CommentID uint   `json:"commentId" binding:"required"`
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

	// Find Comment from CommentID
	var comment models.Comment
	database.DB.First(&comment, json.CommentID)
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

	fmt.Printf("%s has updated a comment.\n\tNew text: %s\n", user.Username, comment.Text)

	// Return updated Comment data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}

/* -------------------------------------------------------------------------- */
/*              DeleteComment | route: comments/delete/:commentId             */
/* -------------------------------------------------------------------------- */
type DeleteCommentRequest struct {
	CommentID uint `uri:"commentId" binding:"required"`
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

	// Find Comment from CommentID
	var comment models.Comment
	database.DB.First(&comment, json.CommentID)
	if comment.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Comment not found."})
		return
	}

	// Find Post from PostID
	var post models.Post
	database.DB.First(&post, comment.PostID)
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
	postUpdatedAt := post.UpdatedAt
	database.DB.Table("comments").Where("post_id = ?", post.ID).Order("created_at DESC, id DESC").First(&lastComment)
	if lastComment.ID != 0 {
		post.CommentedAt = lastComment.CreatedAt
	} else {
		post.CommentedAt = time.Unix(0, 0)
	}
	post.CommentsCount -= 1
	database.DB.Save(&post)

	// By default post.UpdatedAt will get updated by these changes but we don't want that.
	// UpdatedAt should only reflect changes to post.Text
	database.DB.Model(&post).Update("updated_at", postUpdatedAt)

	fmt.Printf("%s has deleted a comment.\n\tPost title: %s\n\tComment text: %s\n", user.Username, post.Title, comment.Text)

	// Return deleted Comment data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}
