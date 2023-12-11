package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type UserController struct {
	UserService       services.UserService
	googleOauthConfig *oauth2.Config
	oauthStateString  string
	ctx               context.Context
}

func NewUserController(userservice services.UserService, googleOauthConfig *oauth2.Config, oauthStateString string) UserController {
	return UserController{
		UserService:       userservice,
		googleOauthConfig: googleOauthConfig,
		oauthStateString:  oauthStateString,
	}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := uc.UserService.CreateUser(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	username := ctx.Param("name")
	uc.UserService.GetUser(&username)
	ctx.JSON(200, "")
}

func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := uc.UserService.UpdateUser(&user); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	username := ctx.Param("name")
	if err := uc.UserService.DeleteUser(&username); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) SignIn(ctx *gin.Context) {
	var signInReq dto.SignInRequest
	if err := ctx.ShouldBindJSON(&signInReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := ctx.Cookie("regular_user_session")
	if err != nil {
		if user == "" {
			uuid := uuid.NewString()
			ctx.SetCookie("regular_user_session", uuid, int(30*time.Minute), "/", "localhost", false, true)
			if err := uc.UserService.SignIn(uuid, signInReq); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
	}
}

func (uc *UserController) SignOut(ctx *gin.Context) {
	uuid, err := ctx.Cookie("regular_user_session")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("LOGIN_ERROR: " + err.Error())})
		return
	}
	ctx.SetCookie("regular_user_session", uuid, -1, "/", "localhost", false, true)
	if err := uc.UserService.SignOut(uuid); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) GoogleSignIn(ctx *gin.Context) {
	url := uc.googleOauthConfig.AuthCodeURL(uc.oauthStateString)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (uc *UserController) GoogleSignInCallback(ctx *gin.Context) {
	// 인증 코드 가져오기
	code := ctx.Query("code")

	// 인증 코드를 사용하여 토큰 교환
	token, err := uc.googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 액세스 토큰을 사용하여 사용자 정보 가져오기
	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	var googleUserInfo models.GoogleUser
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	err = json.NewDecoder(response.Body).Decode(&googleUserInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctx.Cookie("google_user_session")
	if err != nil {
		if user == "" {
			uuid := uuid.NewString()
			ctx.SetCookie("google_user_session", uuid, int(30*time.Minute), "/", "localhost", false, true)
			if err := uc.UserService.GoogleSignIn(uuid, &googleUserInfo); err != nil {
				ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
				return
			}
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
	}
}

func (uc *UserController) GoogleSignOut(ctx *gin.Context) {
	uuid, err := ctx.Cookie("google_user_session")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors.New("GOOGLE_LOGOUT_ERROR: " + err.Error())})
		return
	}
	ctx.SetCookie("google_user_session", uuid, -1, "/", "localhost", false, true)
	if err := uc.UserService.GoogleSignOut(uuid); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	store := cookie.NewStore([]byte(os.Getenv("SECRET")))

	registerroute := rg.Group("/register")
	registerroute.POST("/create", uc.CreateUser)
	registerroute.GET("/get/:name", uc.GetUser)
	registerroute.GET("/getall", uc.GetAll)
	registerroute.PATCH("/update", uc.UpdateUser)
	registerroute.DELETE("/delete/:name", uc.DeleteUser)
	loginroute := rg.Group("login")
	loginroute.Use(sessions.Sessions("mysession", store))
	loginroute.POST("/signin", uc.SignIn)
	loginroute.POST("/signout", uc.SignOut)
	loginroute.GET("/glogin", uc.GoogleSignIn)
	loginroute.GET("/glogincallback", uc.GoogleSignInCallback)
	loginroute.GET("/glogout", uc.GoogleSignOut)
	// privateroute := rg.Group("/private")
	// privateroute.Use(middleware.EnsureLoggedIn())
}
