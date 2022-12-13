package models

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title			string
	Tag				string	
	Text 			string
	Author 			string
	CommentsCount	uint
	CommentedAt		time.Time
	StarsCount		uint
}