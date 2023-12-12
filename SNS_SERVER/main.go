package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers"
	docs "github.com/dik654/Go_projects/SNS_SERVER/docs"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.mongodb.org/mongo-driver/bson"
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
	serviceInstances     services.Services
	controllerInstances  controllers.Controllers
	ctx                  context.Context
	usercollection       *mongo.Collection
	googleusercollection *mongo.Collection
	postcollection       *mongo.Collection
	commentcollection    *mongo.Collection
	likecollection       *mongo.Collection
	mongoclient          *mongo.Client
	redisclient          *redis.Client
	err                  error
)

func init() {
	//
	// ENV bootstrap
	//
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	redisclient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	//
	// GOOGLE_AUTH bootstrap
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
	// MONGO_DB_CLIENT bootstrap
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

	//
	// MONGO_DB_COLLECTION bootstrap
	//

	usercollection = mongoclient.Database("userdb").Collection("users")
	postcollection = mongoclient.Database("postdb").Collection("posts")
	postcollection.Indexes().CreateOne(
		context.Background(),
		mongo.IndexModel{
			Keys: bson.M{
				"title":   "text",
				"content": "text",
				"author":  "text",
			},
		},
	)
	likecollection = mongoclient.Database("likedb").Collection("likes")
	commentcollection = mongoclient.Database("commentdb").Collection("comments")
	googleusercollection = mongoclient.Database("userdb").Collection("google_users")

	//
	// GO_GIN_SERVER bootstrap
	//

	serviceInstances = services.New(redisclient, usercollection, googleusercollection, postcollection, commentcollection, likecollection, ctx)
	controllerInstances = controllers.New(
		serviceInstances,
		googleOauthConfig,
		oauthStateString)
	server = gin.Default()
	server.ForwardedByClientIP = true
	server.SetTrustedProxies([]string{
		"127.0.0.1",
	})
}

// @title			SNS SERVER
// @version		1.0
// @description	mini sns server
// @termsOfService	http://swagger.io/terms/
// @host			localhost:9090
// @BasePath		/v1
func main() {
	defer mongoclient.Disconnect(ctx)

	docs.SwaggerInfo.BasePath = "/v1"
	basepath := server.Group("/v1")
	basepath.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	controllers.RegisterRoutes(controllerInstances, basepath)

	// production 환경에서는 RunTLS로 https 통신을 사용해야함 (쿠키보안 등)
	log.Fatal(server.Run(":9090"))
}
