package models

import "time"

type Like struct {
	ID        uint64    `json:"id" bson:"like_id"`
	UserID    string    `json:"user_id" bson:"like_user_id"`
	UserEmail string    `json:"user_email" bson:"like_user_email"`
	PostID    uint64    `json:"post_id" bson:"like_post_id"`
	CreatedAt time.Time `json:"created_at" bson:"like_created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"like_updated_at"`
}
