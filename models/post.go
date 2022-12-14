package models

import (
	"time"
)

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