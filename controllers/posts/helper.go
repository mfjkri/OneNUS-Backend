package posts

import (
	"math"

	"github.com/mfjkri/OneNUS-Backend/config"
	"github.com/mfjkri/OneNUS-Backend/models"
	"gorm.io/gorm"
)

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

// Fetches posts based on provided configuration
func GetPostsFromContext(dbContext *gorm.DB, perPage uint, pageNumber uint, sortOption string, sortOrder string) ([]models.Post, int64) {
	var posts []models.Post

	// Limit PerPage to config.MAX_PER_PAGE
	clampedPerPage := int64(math.Min(config.MAX_PER_PAGE, float64(perPage)))
	offsetPostsCount := int64(pageNumber-1) * clampedPerPage

	// Get total count for Posts
	var totalPostsCount int64
	dbContext.Count(&totalPostsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetPostsCount < 0) || (offsetPostsCount > totalPostsCount) {
		return posts, 0
	}

	// Sort Posts by sort option provided (defaults to byNew)
	defaultSortOption := config.SORT_BYNEW
	if sortOption == "recent" {
		defaultSortOption = config.SORT_BYRECENT
	} else if sortOption == "hot" {
		defaultSortOption = config.SORT_BYHOT
	}

	// Fetch Posts from [offsetCount, offsetCount + perPage]
	// results order depends on SortOption and SortOrder
	if sortOrder == "ascending" {
		// Reverse page number based on totalPostsCount
		leftOverRecords := math.Min(float64(clampedPerPage), float64(totalPostsCount-offsetPostsCount))
		offsetPostsCount = totalPostsCount - offsetPostsCount - clampedPerPage
		dbContext.Limit(int(leftOverRecords)).Order(defaultSortOption).Offset(int(offsetPostsCount)).Find(&posts)

		// Reverse the page results for descending order
		for i, j := 0, len(posts)-1; i < j; i, j = i+1, j-1 {
			posts[i], posts[j] = posts[j], posts[i]
		}
	} else {
		dbContext.Limit(int(clampedPerPage)).Order(defaultSortOption).Offset(int(offsetPostsCount)).Find(&posts)
	}

	return posts, totalPostsCount
}
