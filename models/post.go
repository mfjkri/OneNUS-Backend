package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title			string
	Tag				string	
	Text 			string
	Author 			string
	RepliesCount	uint
}