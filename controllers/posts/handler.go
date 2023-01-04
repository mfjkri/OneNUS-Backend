package posts

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mfjkri/OneNUS-Backend/config"
	"github.com/mfjkri/OneNUS-Backend/controllers/auth"
	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
	"github.com/mfjkri/OneNUS-Backend/utils"
)

/* -------------------------------------------------------------------------- */
/*                            GetPosts | route: ...                           */
/* -------------------------------------------------------------------------- */
// route: /posts/get/:perPage/:pageNumber/:sortBy/:filterUserId/:filterTag
type GetPostsRequest struct {
	PerPage      uint   `uri:"perPage" binding:"required"`
	PageNumber   uint   `uri:"pageNumber" binding:"required"`
	SortOption   string `uri:"sortOption"`
	SortOrder    string `uri:"sortOrder"`
	FilterUserID uint   `uri:"filterUserId"`
	FilterTag    string `uri:"filterTag"`
}

func GetPosts(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetPostsRequest
	if err := c.ShouldBindUri(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	dbContext := database.DB.Table("posts")

	// Filter database by UserID (if any)
	if json.FilterUserID != 0 {
		targetUser, found := auth.FindUserFromID(c, json.FilterUserID)
		if found == false {
			return
		} else {
			dbContext = dbContext.Where("user_id = ?", targetUser.ID)
		}
	}

	// Filter database by FilterTag (if any)
	if verifyTag(json.FilterTag) {
		dbContext = dbContext.Where("tag = ?", json.FilterTag)
	}

	// Fetch posts
	posts, totalPostsCount := GetPostsFromContext(dbContext, json.PerPage, json.PageNumber, json.SortOption, json.SortOrder)

	// Return fetched posts
	c.JSON(http.StatusAccepted, CreatePostsResponse(&posts, totalPostsCount))
}

/* -------------------------------------------------------------------------- */
/*                GetPostByID | route : /posts/getbyid/:postId                */
/* -------------------------------------------------------------------------- */
type GetPostByIDRequest struct {
	PostID uint `uri:"postId" binding:"required"`
}

func GetPostByID(c *gin.Context) {
	// Check that RequestUser is authenticated
	_, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json GetPostByIDRequest
	if err := c.ShouldBindUri(&json); err != nil {
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

	// Return fetched Post
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

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
	user, found := auth.VerifyAuth(c)
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
	timeNow, canCreatePost := utils.CheckTimeIsAfter(user.LastPostAt, config.USER_POST_COOLDOWN)
	if canCreatePost == false {
		cdLeft := utils.GetCooldownLeft(user.LastPostAt, config.USER_POST_COOLDOWN, timeNow)
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Creating posts too frequently. Please try again in %ds", int(cdLeft.Seconds()))})
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
	post := models.Post{
		Title:         utils.TrimString(strings.TrimSpace(json.Title), config.MAX_POST_TITLE_CHAR),
		Tag:           json.Tag,
		Text:          utils.TrimString(strings.TrimSpace(json.Text), config.MAX_POST_TEXT_CHAR),
		Author:        user.Username,
		User:          user,
		CommentsCount: 0,
		CommentedAt:   time.Unix(0, 0),
		StarsCount:    0,
	}
	new_entry := database.DB.Create(&post)

	// Failed to create entry
	if new_entry.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"message": "Unable to create post. Try again later."})
		return
	}

	// Successfully created a new Post

	// Update PostsCount and LastPostAt for User
	user.PostsCount += 1
	user.LastPostAt = timeNow
	database.DB.Save(&user)

	fmt.Printf("%s has created a post.\n\tPost title: %s\n\tPost text: %s\n", user.Username, post.Title, post.Text)

	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */
/*                  UpdatePostText | route: /posts/updatetext                 */
/* -------------------------------------------------------------------------- */
type UpdatePostTextRequest struct {
	GetPostByIDRequest
	Text string `json:"text" binding:"required"`
}

func UpdatePostText(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := auth.VerifyAuth(c)
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
	timeNow, canCreatePost := utils.CheckTimeIsAfter(user.LastPostAt, config.USER_POST_COOLDOWN)
	if canCreatePost == false {
		cdLeft := utils.GetCooldownLeft(user.LastPostAt, config.USER_POST_COOLDOWN, timeNow)
		c.JSON(http.StatusForbidden, gin.H{"message": fmt.Sprintf("Updating posts too frequently. Please try again in %ds", int(cdLeft.Seconds()))})
		return
	}

	// Find Post from PostID
	var post models.Post
	database.DB.First(&post, json.PostID)
	if post.ID == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "Post not found."})
		return
	}

	// Check User is the author
	if post.UserID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	// Replace Post text and update User LastPostAt
	post.Text = utils.TrimString(strings.TrimSpace(json.Text), config.MAX_POST_TEXT_CHAR)
	user.LastPostAt = timeNow
	database.DB.Save(&post)
	database.DB.Save(&user)

	fmt.Printf("%s has updated a post.\n\tPost title: %s\n\tNew text: %s\n", user.Username, post.Title, post.Text)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}

/* -------------------------------------------------------------------------- */
/*                     DeletePost | route: /delete/:postId                    */
/* -------------------------------------------------------------------------- */
type DeletePostRequest struct {
	PostID uint `uri:"postId" binding:"required"`
}

func DeletePost(c *gin.Context) {
	// Check that RequestUser is authenticated
	user, found := auth.VerifyAuth(c)
	if found == false {
		return
	}

	// Parse RequestBody
	var json DeletePostRequest
	if err := c.ShouldBindUri(&json); err != nil {
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

	// Check User is the author or is admin
	if (post.UserID != user.ID) && (user.Role != config.USER_ROLE_ADMIN) {
		c.JSON(http.StatusForbidden, gin.H{"message": "You do not have valid permissions."})
		return
	}

	database.DB.Delete(&post)

	// Update PostsCount for User
	// user.PostsCount -= 1
	// database.DB.Save(&user)

	fmt.Printf("%s has deleted a post.\n\tPost title: %s\n", user.Username, post.Title)

	// Return new Post data
	c.JSON(http.StatusAccepted, CreatePostResponse(&post))
}
