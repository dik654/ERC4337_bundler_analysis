package dto

type GetPostRequest struct {
	Type uint8  `json:"type" bson:"post_request_type"`
	Body string `json:"body" bson:"post_request_body"`
}
