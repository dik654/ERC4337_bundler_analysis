package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type Services struct {
	UserService    UserService
	PostService    PostService
	CommentService CommentService
}

func New(redisclient *redis.Client, usercollection *mongo.Collection, googleusercollection *mongo.Collection, postcollection *mongo.Collection, commentcollection *mongo.Collection, ctx context.Context) Services {
	userService := NewUserService(redisclient, usercollection, googleusercollection, ctx)
	commentService := NewCommentService(redisclient, commentcollection, ctx)
	postService := NewPostService(commentService, redisclient, postcollection, ctx)
	return Services{
		userService,
		postService,
		commentService,
	}
}
