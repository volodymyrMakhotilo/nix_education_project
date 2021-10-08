package models

type Comment struct {
	PostId uint   `json:"postId"`
	Id     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Body   string `json:"body"`
}
