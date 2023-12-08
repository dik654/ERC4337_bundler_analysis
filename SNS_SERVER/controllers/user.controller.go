package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type UserController struct {
	UserService       services.UserService
	googleOauthConfig *oauth2.Config
	oauthStateString  string
	ctx               context.Context
}

func New(userservice services.UserService, googleOauthConfig *oauth2.Config, oauthStateString string) UserController {
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
	session := sessions.Default(ctx)
	var signInReq dto.SignInRequest
	if err := ctx.ShouldBindJSON(&signInReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.UserService.SignIn(session, signInReq); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) SignOut(ctx *gin.Context) {
	session := sessions.Default(ctx)
	var signInReq dto.SignInRequest
	if err := ctx.ShouldBindJSON(&signInReq); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := uc.UserService.SignOut(session); err != nil {
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
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer response.Body.Close()

	userInfo := struct {
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}{}
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Fprintf(ctx.Writer, "UserInfo: %s\n", userInfo)
}

func (uc *UserController) GoogleSignOut(ctx *gin.Context) {

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
	// privateroute := rg.Group("/private")
	// privateroute.Use(middleware.EnsureLoggedIn())
}
