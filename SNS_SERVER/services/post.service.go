package services

import (
	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
)

type PostService interface {
	CreatePost(*models.Post) error
	GetAllPosts() error
	GetPosts(*dto.GetPostRequest) ([]models.Post, error)
	GetUserPosts() error
	UpdatePost() error
	DeletePost() error
}
