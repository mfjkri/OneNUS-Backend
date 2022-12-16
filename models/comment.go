package models

type Comment struct {
	BaseModel
	
	Text 	string	`json:"text"`

	Author 			string
	User			User 	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	UserID 			uint

	Post			Post	`gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	PostId			uint
}