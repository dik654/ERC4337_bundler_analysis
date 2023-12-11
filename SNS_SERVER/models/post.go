package models

import "time"

type Post struct {
	ID          uint64    `json:"id" bson:"post_id"`
	Title       string    `json:"title" bson:"post_title"`
	Content     string    `json:"content" bson:"post_content"`
	AuthorEmail string    `json:"author" bson:"post_author"`
	AuthorID    string    `json:"author_id" bson:"post_author_id"`
	CreatedAt   time.Time `json:"created_at" bson:"post_created_at"`
	UpdatedAt   time.Time `json:"updated_at" bson:"post_updated_at"`
}
