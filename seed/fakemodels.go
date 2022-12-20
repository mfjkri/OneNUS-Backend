package seed

import (
	"time"

	"github.com/mfjkri/One-NUS-Backend/models"
)

type Post struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Title string
	Tag   string
	Text  string

	Author string
	User   models.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint

	CommentsCount uint
	CommentedAt   time.Time
	StarsCount    uint
}

type Comment struct {
	ID uint `gorm:"primarykey"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Text string `json:"text"`

	Author string
	User   models.User `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID uint

	Post   models.Post `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostId uint
}
