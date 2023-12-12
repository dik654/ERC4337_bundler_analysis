package services

type CommentService interface {
	CreateComment() error
	GetComments() error
	GetUserComments() error
	UpdateComment() error
	DeleteComment() error
}
