package services

type PostService interface {
	CreatePost() error
	GetAllPosts() error
	GetPosts() error
	GetUserPosts() error
	UpdatePost() error
	DeletePost() error
}
