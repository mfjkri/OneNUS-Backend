package controllers

import (
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
func verifyTag(tag string) (valid bool) {
	valid = false
	for _, x := range models.ValidTags {
		if x == tag {
			valid = true
		}
	}
	return
}

type PostResponse struct {
	ID            uint   `json:"id" binding:"required"`
	Title         string `json:"title" binding:"required"`
	Tag           string `json:"tag" binding:"required"`
	Text          string `json:"text" binding:"required"`
	Author        string `json:"author" binding:"required"`
	UserID        uint   `json:"userId" binding:"required"`
	CommentsCount uint   `json:"commentsCount" binding:"required"`
	CommentedAt   int64  `json:"commentedAt" binding:"required"`
	StarsCount    uint   `json:"starsCount" binding:"required"`
	CreatedAt     int64  `json:"createdAt" binding:"required"`
	UpdatedAt     int64  `json:"updatedAt" binding:"required"`
}

// Convert a Post Model into a JSON format
func CreatePostResponse(post *models.Post) PostResponse {
	return PostResponse{
		ID:            post.ID,
		Title:         post.Title,
		Tag:           post.Tag,
		Text:          post.Text,
		Author:        post.Author,
		UserID:        post.UserID,
		CommentsCount: post.CommentsCount,
		CommentedAt:   post.CommentedAt.Unix(),
		StarsCount:    post.StarsCount,
		CreatedAt:     post.CreatedAt.Unix(),
		UpdatedAt:     post.UpdatedAt.Unix(),
	}
}

type GetPostsResponse struct {
	Posts      []PostResponse `json:"posts" binding:"required"`
	PostsCount int64          `json:"postsCount" binding:"required"`
}

// Bundles and convert multiple Post models into a JSON format
func CreatePostsResponse(posts *[]models.Post, totalPostsCount int64) GetPostsResponse {
	var postsResponse []PostResponse
	for _, post := range *posts {
		postReponse := CreatePostResponse(&post)
		postsResponse = append(postsResponse, postReponse)
	}

	return GetPostsResponse{
		Posts:      postsResponse,
		PostsCount: totalPostsCount,
	}
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*    GetPosts | route: /posts/get/:perPage/:pageNumber/:sortBy/:filterTag    */
/* -------------------------------------------------------------------------- */
type GetPostsRequest struct {
	PerPage    uint   `uri:"perPage" binding:"required"`
	PageNumber uint   `uri:"pageNumber" binding:"required"`
	SortBy     string `uri:"sortBy"`
	FilterTag  string `uri:"filterTag"`
}

func GetPosts(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetPostsRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Limit PerPage to 10
	perPage := int64(math.Min(MAX_PER_PAGE, float64(json.PerPage)))
	offsetPostCount := int64(json.PageNumber-1) * perPage

	dbContext := database.DB.Table("posts")

	// Filter database by FilterTag (if any)
	if verifyTag(json.FilterTag) {
		dbContext = dbContext.Where("tag = ?", json.FilterTag)
	}

	// Get total count for Posts
	var totalPostsCount int64
	dbContext.Count(&totalPostsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetPostCount < 0) || (offsetPostCount > totalPostsCount) {
		c.JSON(http.StatusForbidden, gin.H{"message": "No more posts found."})
		return
	}

	// Sort Posts by sort option provided (defaults to byNew)
	defaultSortOption := ByNew
	if json.SortBy == "byRecent" {
		defaultSortOption = ByRecent
	} else if json.SortBy == "byHot" {
		defaultSortOption = ByHot
	}

	// Fetch Posts from [offsetCount, offsetCount + perPage]
	var posts []models.Post
	dbContext.Limit(int(perPage)).Order(defaultSortOption).Offset(int(offsetPostCount)).Find(&posts)

	// Return fetched posts
	c.JSON(http.StatusAccepted, CreatePostsResponse(&posts, totalPostsCount))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                GetPostByID | route : /posts/getbyid/:postId                */
/* -------------------------------------------------------------------------- */
type GetPostByIDRequest struct {
	PostId uint `uri:"postId" binding:"required"`
}

func GetPostByID(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetPostByIDRequest
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

	// Return fetched Post
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                        CreatePost | route: /post/get                       */
/* -------------------------------------------------------------------------- */
type CreatePostRequest struct {
	Title string `json:"title" binding:"required"`
	Tag   string `json:"tag" binding:"required"`
	Text  string `json:"text" binding:"required"`
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
	timeNow, canCreatePost := utils.CheckTimeIsAfter(user.LastPostAt, USER_POST_COOLDOWN)
	if canCreatePost == false {
		c.JSON(http.StatusForbidden, gin.H{"message": "Creating posts too frequently. Please try again later."})
		return
	}

	// Check that Title and Text does not contain illegal characters
	if !(utils.ContainsValidCharactersOnly(json.Title) && utils.ContainsValidCharactersOnly(json.Text)) {
		c.JSON(http.StatusForbidden, gin.H{"message": "Title or Body contains illegal characters."})
		return
	}

	// Check that the Tag provided is valid
	validTag := verifyTag(json.Tag)
	if validTag == false {
		c.JSON(http.StatusForbidden, gin.H{"message": "Unknown tag for post."})
		return
	}

	// Try to create new Post
	initialCommentsCount := uint(0)
	initialStarsCount := uint(0)
	post := models.Post{
		Title:         utils.TrimString(json.Title, MAX_POST_TITLE_CHAR),
		Tag:           json.Tag,
		Text:          utils.TrimString(json.Text, MAX_POST_TEXT_CHAR),
		Author:        user.Username,
		User:          user,
		CommentsCount: initialCommentsCount,
		CommentedAt:   time.Unix(0, 0),
		StarsCount:    initialStarsCount,
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

	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                  UpdatePostText | route: /posts/updatetext                 */
/* -------------------------------------------------------------------------- */
type UpdatePostTextRequest struct {
	GetPostByIDRequest
	Text string `json:"text" binding:"required"`
}

func UpdatePostText(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json UpdatePostTextRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	// Prevent frequent UpdatePostText by User
	timeNow, canCreatePost := utils.CheckTimeIsAfter(user.LastPostAt, USER_POST_COOLDOWN)
	if canCreatePost == false {
		c.JSON(http.StatusForbidden, gin.H{"message": "Updating posts too frequently. Please try again later."})
		return
	}

	// Find Post from PostId
	var post models.Post
	database.DB.First(&post, json.PostId)
	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	// Check User is the author
	if strings.ToLower(post.Author) != user.Username {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Replace Post text and update User LastPostAt
	post.Text = utils.TrimString(json.Text, MAX_POST_TEXT_CHAR)
	user.LastPostAt = timeNow
	database.DB.Save(&post)
	database.DB.Save(&user)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */

/* -------------------------------------------------------------------------- */
/*                     DeletePost | route: /delete/:postId                    */
/* -------------------------------------------------------------------------- */
type DeletePostRequest struct {
	PostId uint `uri:"postId" binding:"required"`
}

func DeletePost(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json DeletePostRequest
	if err := c.ShouldBindUri(&json); err != nil {
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

	// Check User is the author
	if strings.ToLower(post.Author) != user.Username {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	database.DB.Delete(&post)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */
