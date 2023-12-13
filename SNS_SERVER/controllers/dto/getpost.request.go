package dto

type GetPostRequest struct {
	Type uint8  `json:"type"`
	Body string `json:"body"`
}
