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
	LikeService    LikeService
}

func New(redisclient *redis.Client, usercollection *mongo.Collection, googleusercollection *mongo.Collection, postcollection *mongo.Collection, commentcollection *mongo.Collection, likecollection *mongo.Collection, ctx context.Context) Services {
	userService := NewUserService(redisclient, usercollection, googleusercollection, ctx)
	postService := NewPostService(redisclient, postcollection, ctx)
	commentService := NewCommentService(redisclient, commentcollection, ctx)
	likeService := NewLikeService(redisclient, likecollection, ctx)
	return Services{
		userService,
		postService,
		commentService,
		likeService,
	}
}
