package dto

type PostLikeRequest struct {
	PostID string `json:"id"`
	Like   bool   `json:"like"`
}
