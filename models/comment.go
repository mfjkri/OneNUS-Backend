package models

import (
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	BaseModel

	Text string `json:"text"`

	Author string
	User   User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint

	Post   Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostID uint
}

func (comment *Comment) AfterDelete(tx *gorm.DB) (err error) {
	var post Post
	tx.First(&post, comment.PostID)
	if post.ID != 0 {
		// Update associated Post metadata to reflect new changes
		postUpdatedAt := post.UpdatedAt

		post.CommentsCount -= 1

		var lastComment Comment
		tx.Table("comments").Where("post_id = ?", post.ID).Order("created_at DESC, id DESC").First(&lastComment)
		if lastComment.ID != 0 {
			post.CommentedAt = lastComment.CreatedAt
		} else {
			post.CommentedAt = time.Unix(0, 0)
		}

		tx.Save(&post)

		// By default post.UpdatedAt will get updated by these changes but we don't want that.
		// UpdatedAt should only reflect changes to post.Text
		tx.Model(&post).Update("updated_at", postUpdatedAt)
	}

	// var user User
	// tx.First(&user, comment.UserID)
	// if user.ID != 0 {
	// 	// Decrement user comments count
	// 	user.CommentsCount -= 1
	// 	tx.Save(&user)
	// }

	return
}
