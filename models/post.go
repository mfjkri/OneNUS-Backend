package models

import (
	"time"
)

var ValidTags = [4]string{"general", "cs", "life", "misc"}

type Post struct {
	BaseModel
	
	Title			string
	Tag				string	
	Text 			string
	
	Author 			string
	User			User 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID 			uint

	CommentsCount	uint
	CommentedAt		time.Time
	StarsCount		uint
}