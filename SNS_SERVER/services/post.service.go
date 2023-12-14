package services

import (
	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
)

type PostService interface {
	CreatePost(*dto.CreatePostRequest, *dto.SessionInfo) error
	GetAllPosts(*dto.PaginationRequest) ([]models.Post, error)
	GetPosts(*dto.GetPostRequest) ([]models.Post, error)
	UpdatePost(models.Post) error
	DeletePost(string) error
	LikePost(*dto.PostLikeRequest) error
	JudgePost(models.Judge) error
	CanEditPost(*dto.SessionInfo, string) (bool, error)
}
