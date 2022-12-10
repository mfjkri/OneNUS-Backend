package models

type User struct {
	Id       uint   `json:"id"`
	Username     string `json:"name"`
	Password []byte `json:"-"`
}