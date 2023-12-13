package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type CommentServiceImpl struct {
	redisclient       *redis.Client
	commentcollection *mongo.Collection
	ctx               context.Context
}

func NewCommentService(redisclient *redis.Client, commentcollection *mongo.Collection, ctx context.Context) CommentService {
	return &CommentServiceImpl{
		redisclient:       redisclient,
		commentcollection: commentcollection,
		ctx:               ctx,
	}
}

func (c *CommentServiceImpl) CreateComment() error {
	return nil
}

func (c *CommentServiceImpl) GetComments() error {
	return nil
}

func (c *CommentServiceImpl) UpdateComment() error {
	return nil
}

func (c *CommentServiceImpl) DeleteComment() error {
	return nil
}
