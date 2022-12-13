package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username 		string		`gorm:"unique"`
	Password 		[]byte 	
	LastPostAt		time.Time 	`gorm:"autoCreateTime"` 
	LastCommentAt	time.Time 	`gorm:"autoCreateTime"` 
}