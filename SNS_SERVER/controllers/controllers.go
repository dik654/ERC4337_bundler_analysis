package controllers

import (
	"errors"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type Controllers struct {
	UserController    UserController
	PostController    PostController
	CommentController CommentController
}

// 세션 정보를 추출하는 함수
func GetSessionInfo(ctx *gin.Context) (*dto.SessionInfo, error) {
	// 세션 타입과 키를 매핑
	sessionMap := map[string]string{
		"regular_user_session": "regular",
		"google_user_session":  "google",
	}

	for key, userType := range sessionMap {
		if sessionUUID, err := ctx.Cookie(key); err == nil {
			return &dto.SessionInfo{SessionUUID: sessionUUID, UserType: userType}, nil
		}
	}
	return nil, errors.New("GET_SESSION_INFO_ERROR: no active session")
}

func New(redisclient *redis.Client, services services.Services, googleOauthConfig *oauth2.Config, oauthStateString string) Controllers {
	userController := NewUserController(redisclient, services.UserService, googleOauthConfig, oauthStateString)
	postController := NewPostController(services.PostService)
	commentController := NewCommentController(services.CommentService)
	return Controllers{
		userController,
		postController,
		commentController,
	}
}

func RegisterRoutes(controllers Controllers, rg *gin.RouterGroup) {
	controllers.UserController.RegisterUserRoutes(rg)
	controllers.PostController.RegisterPostRoutes(rg)
	controllers.CommentController.RegisterCommentRoutes(rg)
}
