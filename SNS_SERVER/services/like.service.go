package services

type LikeService interface {
	LikePost() error
	UnLikePost() error
	GetLikes() error
}
