package models

import "gorm.io/gorm"

type Post struct {
	gorm.Model
	Title	string	`json:"title"`
	Tag		string 	`json:"tag"`
	Text 	string	`json:"text"`
}