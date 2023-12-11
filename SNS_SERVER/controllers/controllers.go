package controllers

import (
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Controllers struct {
	UserController    UserController
	PostController    PostController
	CommentController CommentController
	LikeController    LikeController
}

func New(services services.Services, googleOauthConfig *oauth2.Config, oauthStateString string) Controllers {
	userController := NewUserController(services.UserService, googleOauthConfig, oauthStateString)
	postController := NewPostController(services.PostService)
	commentController := NewCommentController(services.CommentService)
	likeController := NewLikeController(services.LikeService)
	return Controllers{
		userController,
		postController,
		commentController,
		likeController,
	}
}

func RegisterRoutes(controllers Controllers, rg *gin.RouterGroup) {
	controllers.UserController.RegisterUserRoutes(rg)
	controllers.PostController.RegisterPostRoutes(rg)
	controllers.CommentController.RegisterCommentRoutes(rg)
	controllers.LikeController.RegisterLikeRoutes(rg)
}
