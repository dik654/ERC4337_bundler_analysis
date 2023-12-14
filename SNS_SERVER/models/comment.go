package models

import "time"

type Comment struct {
	ID        string    `json:"id" bson:"comment_id"`
	AuthorID  string    `json:"author_id" bson:"comment_author_id"`
	PostID    string    `json:"post_id" bson:"comment_post_id"`
	Content   string    `json:"content" bson:"comment_content"`
	CreatedAt time.Time `json:"created_at" bson:"comment_created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"comment_updated_at"`
}
