package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	server               *gin.Engine
	googleOauthConfig    *oauth2.Config
	oauthStateString     string
	userservice          services.UserService
	usercontroller       controllers.UserController
	ctx                  context.Context
	usercollection       *mongo.Collection
	googleusercollection *mongo.Collection
	mongoclient          *mongo.Client
	err                  error
)

func init() {

	//
	// GODOTENV
	//

	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	//
	// GOOGLE_OAUTH
	//

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:9090/v1/login/glogincallback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	oauthStateString = os.Getenv("SECRET_KEY")

	//
	// REDIS
	//

	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))

	//
	// MONGO_DB
	//

	ctx = context.Background()

	mongoconn := options.Client().ApplyURI("mongodb://localhost:27017")
	mongoclient, err = mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal(err)
	}
	err = mongoclient.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("mongo connection established")

	usercollection = mongoclient.Database("userdb").Collection("users")
	googleusercollection = mongoclient.Database("userdb").Collection("google_users")
	usercontroller = controllers.New(userservice, store, googleOauthConfig, oauthStateString)
	userservice = services.NewUserService(usercollection, googleusercollection, ctx)

	//
	// GO GIN
	//

	server = gin.Default()
	server.ForwardedByClientIP = true
	server.SetTrustedProxies([]string{
		"127.0.0.1",
	})
}

func main() {
	defer mongoclient.Disconnect(ctx)

	basepath := server.Group("/v1")
	usercontroller.RegisterUserRoutes(basepath)

	// production 환경에서는 RunTLS로 https 통신을 사용해야함 (쿠키보안 등)
	log.Fatal(server.Run(":9090"))
}
