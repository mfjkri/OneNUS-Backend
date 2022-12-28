package models

import (
	"time"
)

type User struct {
	BaseModel

	Username string `gorm:"unique"`
	Password []byte
	Role     string
	Bio      string
	Private  bool `gorm:"default:false"`

	PostsCount    uint `gorm:"default:0"`
	CommentsCount uint `gorm:"default:0"`

	LastPostAt    time.Time `gorm:"autoCreateTime"`
	LastCommentAt time.Time `gorm:"autoCreateTime"`
}
