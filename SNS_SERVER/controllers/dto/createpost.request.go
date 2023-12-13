package dto

type CreatePostRequest struct {
	Title   string `json:"title" bson:"post_createrequest_title"`
	Content string `json:"content" bson:"post_createrequest_content"`
}
