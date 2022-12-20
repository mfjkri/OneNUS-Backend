package models

import (
	"time"
)

type User struct {
	BaseModel

	Username string `gorm:"unique"`
	Password []byte
	Role     string

	LastPostAt    time.Time `gorm:"autoCreateTime"`
	LastCommentAt time.Time `gorm:"autoCreateTime"`
}
