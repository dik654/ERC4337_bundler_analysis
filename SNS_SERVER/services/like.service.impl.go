package services

import (
	"context"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

type LikeServiceImpl struct {
	redisclient    *redis.Client
	likecollection *mongo.Collection
	ctx            context.Context
}

func NewLikeService(redisclient *redis.Client, likecollection *mongo.Collection, ctx context.Context) LikeService {
	return &LikeServiceImpl{
		redisclient:    redisclient,
		likecollection: likecollection,
		ctx:            ctx,
	}
}

func (l *LikeServiceImpl) LikePost() error {
	return nil
}

func (l *LikeServiceImpl) UnLikePost() error {
	return nil
}

func (l *LikeServiceImpl) GetLikes() error {
	return nil
}
