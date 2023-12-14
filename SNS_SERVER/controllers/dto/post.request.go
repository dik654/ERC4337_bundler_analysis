package dto

type CreatePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type GetPostRequest struct {
	Type uint8  `json:"type"`
	Body string `json:"body"`
}
