package seed

import (
	"fmt"
	"time"

	"github.com/mfjkri/OneNUS-Backend/database"
	"github.com/mfjkri/OneNUS-Backend/models"
)

func UpdateUsers() {
	fmt.Println("Updating users data...")
	var users []models.User
	database.DB.Find(&users)

	for _, user := range users {
		var totalPostsCount int64
		database.DB.Table("posts").Where("user_id = ?", user.ID).Count(&totalPostsCount)
		user.PostsCount = uint(totalPostsCount)

		var totalCommentsCount int64
		database.DB.Table("comments").Where("user_id = ?", user.ID).Count(&totalCommentsCount)
		user.CommentsCount = uint(totalCommentsCount)

		database.DB.Save(&user)
	}
	fmt.Println("Update users data complete!")
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
		} else {
			post.CommentedAt = time.Unix(0, 0)
		}

		postUpdatedAt := post.UpdatedAt
		database.DB.Save(&post)
		database.DB.Model(&post).Update("updated_at", postUpdatedAt)
	}

	fmt.Println("Update posts data complete!")
}
