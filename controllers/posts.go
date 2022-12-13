package controllers

import (
	"math"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
	"github.com/mfjkri/One-NUS-Backend/utils"
)

/* -------------------------------------------------------------------------- */
/*                              Helper functions                              */
/* -------------------------------------------------------------------------- */
var ValidTags = [4]string{"general", "cs", "life", "misc"}

func verifyTag(tag string) (valid bool) {
	valid = false
	for _, x := range ValidTags {
		if x == tag {
			valid = true
		}
	}
	return
}
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                        CreatePost | route: /post/get                       */
/* -------------------------------------------------------------------------- */
type CreatePostRequest struct {
	Title			string	`json:"title" binding:"required"`
	Tag				string 	`json:"tag" binding:"required"`
	Text 			string	`json:"text" binding:"required"`
}

func CreatePost(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody 
	var json CreatePostRequest
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
      return
    }

	// Prevent frequent CreatePosts by User
	timeNow, canCreatePost := utils.CheckTimeIsAfter(user.LastPostAt, 1 * time.Minute)
	if canCreatePost == false {
		c.JSON(http.StatusForbidden, gin.H{"message": "Creating posts too frequently. Please try again later."})
      	return
	}

	// Check that the Tag provided is valid
	validTag := verifyTag(json.Tag)
	if validTag == false {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unknown tag for post."})
      	return
	}

	// Try to create new Post
	initialRepliesCount := uint(0) 
	post := models.Post{
		Title: json.Title,
		Tag: json.Tag,
		Text: json.Text,
		Author: user.Username,
		RepliesCount: initialRepliesCount,
	}
	new_entry := database.DB.Create(&post)

	// Failed to create entry
	if new_entry.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Unable to create post. Try again later."})
      	return
	}

	// Successfully created a new Post

	// Update LastPostAt for User
	user.LastPostAt = timeNow
	database.DB.Save(&user)

	c.JSON(http.StatusAccepted, gin.H{
		"id": user.ID,
		"username": user.Username,
	})
}
/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                        GetPosts | route: /posts/get                        */
/* -------------------------------------------------------------------------- */
type GetPostsRequest struct {
	PerPage		uint 	`form:"perPage" json:"perPage" binding:"required"`
	PageNumber 	uint 	`form:"pageNumber" json:"pageNumber" binding:"required"`
}

func GetPosts(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody 
	var json GetPostsRequest
    if err := c.ShouldBindJSON(&json); err != nil {
      c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
      return
    }

	// Limit PerPage to 10
	perPage := int64(math.Min(10, float64(json.PageNumber)))
	offsetPostCount := int64(json.PageNumber - 1) * perPage

	// Get total count for Posts
	var totalPostsCount int64
	database.DB.Model(&models.Post{}).Count(&totalPostsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetPostCount < 0) || (offsetPostCount > totalPostsCount - perPage) {
		c.JSON(http.StatusForbidden, gin.H{"message": "No more posts found."})
		return
	}

	// Fetch Posts from [offsetCount, offsetCount + perPage]
	var posts []models.Post
	database.DB.Limit(int(perPage)).Offset(int(offsetPostCount)).Find(&posts)

	// Return fetched posts
	c.JSON(http.StatusAccepted, gin.H{
		"id": user.ID,
		"username": user.Username,
	})
}
/* -------------------------------------------------------------------------- */