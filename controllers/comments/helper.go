package comments

import (
	"math"

	"github.com/mfjkri/OneNUS-Backend/config"
	"github.com/mfjkri/OneNUS-Backend/models"
	"gorm.io/gorm"
)

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

	// Limit PerPage to config.MAX_PER_PAGE
	clampedPerPage := int64(math.Min(config.MAX_PER_PAGE, float64(perPage)))
	offsetCommentsCount := int64(pageNumber-1) * clampedPerPage

	// Get total count for Comments
	var totalCommentsCount int64
	dbContext.Count(&totalCommentsCount)

	// If we are request beyond the bounds of total count, error
	if (offsetCommentsCount < 0) || (offsetCommentsCount > totalCommentsCount) {
		return comments, 0
	}

	// Sort Comments by sort option provided (defaults to byNew)
	defaultSortOption := config.SORT_BYNEW
	if sortOption == "recent" {
		defaultSortOption = config.SORT_BYRECENT
	}

	// Fetch Comments from [offsetCount, offsetCount + perPage]
	// results order depends on SortOption and SortOrder
	if sortOrder == "ascending" {
		// Reverse page number based on totalPostsCount
		leftOverRecords := math.Min(float64(clampedPerPage), float64(totalCommentsCount-offsetCommentsCount))
		offsetCommentsCount = totalCommentsCount - offsetCommentsCount - clampedPerPage
		dbContext.Limit(int(leftOverRecords)).Order(defaultSortOption).Offset(int(offsetCommentsCount)).Find(&comments)

		// Reverse the page results for descending order
		for i, j := 0, len(comments)-1; i < j; i, j = i+1, j-1 {
			comments[i], comments[j] = comments[j], comments[i]
		}
	} else {
		dbContext.Limit(int(clampedPerPage)).Order(defaultSortOption).Offset(int(offsetCommentsCount)).Find(&comments)
	}

	return comments, totalCommentsCount
}
