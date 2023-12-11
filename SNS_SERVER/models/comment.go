package models

import "time"

type Comment struct {
	ID        uint64    `json:"id" bson:"comment_id"`
	UserID    string    `json:"user_id" bson:"comment_user_id"`
	UserEmail string    `json:"user" bson:"comment_user_email"`
	PostID    uint64    `json:"post_id" bson:"comment_post_id"`
	Body      string    `json:"body" bson:"comment_body"`
	CreatedAt time.Time `json:"created_at" bson:"comment_created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"comment_updated_at"`
}
