package comments

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/config"
	"github.com/mfjkri/OneNUS-Backend/controllers/auth"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"github.com/mfjkri/OneNUS-Backend/utils"
)

/* -------------------------------------------------------------------------- */
/*                          GetComments | route: ...                          */
/* -------------------------------------------------------------------------- */
// route: /comments/get/:postId/:perPage/:pageNumber/:sortOption/:sortOrder
type GetCommentsRequest struct {
	PostID     uint   `uri:"postId" binding:"required"`
	PerPage    uint   `uri:"perPage" binding:"required"`
	PageNumber uint   `uri:"pageNumber" binding:"required"`
	SortOption string `uri:"sortOption"`
	SortOrder  string `uri:"sortOrder"`
}

func GetComments(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := auth.VerifyAuth(c)
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
		c.JSON(http.StatusForbidden, gin.H{"message": "Post not found."})
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
	user, found := auth.VerifyAuth(c)
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
	timeNow, canCreateComment := utils.CheckTimeIsAfter(user.LastCommentAt, config.USER_COMMENT_COOLDOWN)
	if canCreateComment == false {
		cdLeft := utils.GetCooldownLeft(user.LastCommentAt, config.USER_COMMENT_COOLDOWN, timeNow)
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
		Text:   utils.TrimString(strings.TrimSpace(json.Text), config.MAX_COMMENT_TEXT_CHAR),
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

	// Update CommentsCount and LastCommentAt for User
	user.CommentsCount += 1
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
	user, found := auth.VerifyAuth(c)
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
	timeNow, canUpdateComment := utils.CheckTimeIsAfter(user.LastCommentAt, config.USER_COMMENT_COOLDOWN)
	if canUpdateComment == false {
		cdLeft := utils.GetCooldownLeft(user.LastCommentAt, config.USER_COMMENT_COOLDOWN, timeNow)
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Updating comments too frequently. Please try again in %ds", int(cdLeft.Seconds()))})
		return
	}

	// Check User is the author
	if comment.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Replace Comment text and update User LastCommentAt
	comment.Text = utils.TrimString(strings.TrimSpace(json.Text), config.MAX_COMMENT_TEXT_CHAR)
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
	user, found := auth.VerifyAuth(c)
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

	// Check User is the author or is admin
	if (comment.UserID != user.ID) && (user.Role != config.USER_ROLE_ADMIN) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Delete Comment
	database.DB.Delete(&comment)

	fmt.Printf("%s has deleted a comment.\n\tComment text: %s\n", user.Username, comment.Text)

	// Return deleted Comment data
	c.JSON(http.StatusAccepted, CreateCommentResponse(&comment))
}
