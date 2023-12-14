package dto

type CreateCommentRequest struct {
	PostID  string `json:"id"`
	Content string `json:"content"`
}

type GetCommentRequest struct {
	Type uint8  `json:"type"`
	Body string `json:"body"`
}
