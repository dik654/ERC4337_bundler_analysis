package services

import (
	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
)

type UserService interface {
	CreateUser(*models.User) error
	GetUser(*string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(*models.User) error
	DeleteUser(*string) error
	SignIn(string, dto.SignInRequest) error
	SignOut(string) error
	GoogleSignIn(string, *models.GoogleUser) error
	GoogleSignOut(string) error
}
