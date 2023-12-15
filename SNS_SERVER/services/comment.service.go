package services

import (
	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
)

type CommentService interface {
	CreateComment(*dto.CreateCommentRequest, *dto.SessionInfo) error
	GetComments(*dto.GetCommentRequest) ([]models.Comment, error)
	UpdateComment(*dto.CreateCommentRequest) error
	DeleteComment(string) error
	DeleteComments(string) error
	CanEditPost(*dto.SessionInfo, string) (bool, error)
}
