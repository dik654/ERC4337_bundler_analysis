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
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"
)

type UserController struct {
	redisClient       *redis.Client
	UserService       services.UserService
	googleOauthConfig *oauth2.Config
	oauthStateString  string
	ctx               context.Context
}

func NewUserController(redisclient *redis.Client, userservice services.UserService, googleOauthConfig *oauth2.Config, oauthStateString string) UserController {
	return UserController{
		redisClient:       redisclient,
		UserService:       userservice,
		googleOauthConfig: googleOauthConfig,
		oauthStateString:  oauthStateString,
	}
}

// RegisterUser godoc
//
//	@Summary		sign up regular user
//	@Tags			register regular user
//	@Description	write user informations to mongodb
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User Data"
//	@Success		200		{string}	success
//	@Router			/register/create [post]
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

// GetUser godoc
//
//	@Summary		get user data
//	@Tags			register regular user
//	@Description	get regular user informations to mongodb
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"regular_user_id"
//	@Success		200	{object}	models.User
//	@Router			/register/get/{id} [get]
func (uc *UserController) GetUser(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := uc.UserService.GetUser(&id)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

// GetAllUser godoc
//
//	@Summary		get user data
//	@Tags			register regular user
//	@Description	get all regular user informations to mongodb
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	[]models.User
//	@Router			/register/getall [get]
func (uc *UserController) GetAll(ctx *gin.Context) {
	users, err := uc.UserService.GetAll()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

// UpdateUser godoc
//
//	@Summary		update regular user data
//	@Tags			register regular user
//	@Description	update regular user informations to mongodb
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User data"
//	@Success		200		{string}	success
//	@Router			/register/update [patch]
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

// DeleteUser godoc
//
//	@Summary		delete regular user data
//	@Tags			register regular user
//	@Description	delete user information in mongodb
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"regular_user_id"
//	@Success		200	{string}	success
//	@Router			/register/delete/{id} [delete]
func (uc *UserController) DeleteUser(ctx *gin.Context) {
	id := ctx.Param("id")
	if err := uc.UserService.DeleteUser(&id); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SignIn godoc
//
//	@Summary		sign in regular user
//	@Tags			sign in/out regular user
//	@Description	sign in regular user
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.SignInRequest	true	"regular login data"
//	@Success		200		{string}	success
//	@Header			200		{string}	Set-Cookie	"Session Cookie"
//	@Router			/login/signin [post]
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
			return
		} else {
			ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
			return
		}
	} else {
		_, err = uc.redisClient.Get(ctx, "regular_user_session:"+user).Result()
		if err != nil {
			ctx.SetCookie("regular_user_session", user, -1, "/", "localhost", false, true)
			ctx.JSON(http.StatusBadGateway, gin.H{"message": "세션이 만료되었습니다. " + err.Error()})
			return
		}
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// SignOut godoc
//
//	@Summary		sign out regular user
//	@Tags			sign in/out regular user
//	@Description	sign out regular user
//	@Accept			json
//	@Produce		json
//	@Param			regular_user_session	header		string	true	"Session Cookie"
//	@Success		200						{string}	success
//	@Router			/login/signout [post]
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

// GoogleSignIn godoc
//
//	@Summary		sign up google user
//	@Tags			sign in/out google user
//	@Description	sign out google user
//	@Accept			json
//	@Produce		json
//	@Success		200	{string}	success
//	@Header			200	{string}	Set-Cookie	"Session Cookie"
//	@Router			/login/glogin [get]
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
	} else {
		_, err = uc.redisClient.Get(ctx, "google_user_session:"+user).Result()
		if err != nil {
			ctx.SetCookie("google_user_session", user, -1, "/", "localhost", false, true)
			ctx.JSON(http.StatusBadGateway, gin.H{"message": "세션이 만료되었습니다. " + err.Error()})
		}
	}
}

// GoogleSignOut godoc
//
//	@Summary		sign out google user
//	@Tags			sign in/out google user
//	@Description	sign out google user
//	@Accept			json
//	@Produce		json
//	@Param			google_user_session	header		string	true	"Session Cookie"
//	@Success		200					{string}	success
//	@Router			/login/glogout [get]
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
	registerroute.GET("/get/:id", uc.GetUser)
	registerroute.GET("/getall", uc.GetAll)
	registerroute.PATCH("/update", uc.UpdateUser)
	registerroute.DELETE("/delete/:id", uc.DeleteUser)
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
