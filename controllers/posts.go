package controllers

import (
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
)


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

type CreatePostRequest struct {
	Title			string	`json:"title" binding:"required"`
	Tag				string 	`json:"tag" binding:"required"`
	Text 			string	`json:"text" binding:"required"`
}

// CreatePost | route:/post/create
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
	c.JSON(http.StatusAccepted, gin.H{
		"id": user.ID,
		"username": user.Username,
	})
}


type GetPostsRequest struct {
	PerPage		uint 	`form:"perPage" json:"perPage" binding:"required"`
	PageNumber 	uint 	`form:"pageNumber" json:"pageNumber" binding:"required"`
}

// GetPosts | route:/posts/get
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