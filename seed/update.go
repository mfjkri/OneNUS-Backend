package seed

import (
	"fmt"

	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
)

func UpdateUsers() {
	fmt.Println("Updating users data...")
	var users []models.User
	database.DB.Find(&users)

	for _, user := range users {
		if user.Username != "admin" {
			user.Role = "member"
		} else {
			user.Role = "admin"
		}

		database.DB.Save(&user)
	}
	fmt.Println("Update complete!")
}

func UpdatePosts() {
	fmt.Println("Updating posts data...")

	var posts []models.Post
	database.DB.Find(&posts)

	for _, post := range posts {
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
