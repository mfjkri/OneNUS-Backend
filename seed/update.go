package seed

import (
	"fmt"

	"github.com/mfjkri/One-NUS-Backend/database"
	"github.com/mfjkri/One-NUS-Backend/models"
)

func UpdateData() {
	fmt.Println("Updating posts data...")

	var posts []models.Post
	database.DB.Find(&posts)

	for _, post := range(posts) {
		dbContext := database.DB.Table("comments").Where("post_id = ?", post.ID)

		var totalCommentsCount int64
		dbContext.Count(&totalCommentsCount) 
		post.CommentsCount = uint(totalCommentsCount)
		
		var lastComment models.Comment
		dbContext.Order("created_at DESC, id DESC").First(&lastComment)
		if lastComment.ID != 0 {
			post.CommentedAt = lastComment.CreatedAt
		}
		
		
		database.DB.Save(&post)
		database.DB.Model(&post).Update("updated_at", post.CreatedAt)
	}

	fmt.Println("Update complete!")
}