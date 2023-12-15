package models

import "time"

type Judge struct {
	Type   uint8  `json:"type"`
	PostID string `json:"id"`
	Body   string `json:"body"`
}

type Post struct {
	ID        string    `json:"id" bson:"post_id"`
	AuthorID  string    `json:"author_id" bson:"post_author_id"`
	Title     string    `json:"title" bson:"post_title"`
	Content   string    `json:"content" bson:"post_content"`
	Like      uint64    `json:"like" bson:"post_like"`
	UnLike    uint64    `json:"unlike" bson:"post_unlike"`
	Judge     []Judge   `json:"judge" bson:"post_judge"`
	CreatedAt time.Time `json:"created_at" bson:"post_created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"post_updated_at"`
}
