package services

type CommentService interface {
	CreateComment() error
	GetComments() error
	UpdateComment() error
	DeleteComment() error
}
