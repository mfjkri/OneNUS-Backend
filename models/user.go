package models

import (
	"time"
)

type User struct {
	BaseModel

	Username      string `gorm:"unique"`
	Password      []byte
	LastPostAt    time.Time `gorm:"autoCreateTime"`
	LastCommentAt time.Time `gorm:"autoCreateTime"`
}
