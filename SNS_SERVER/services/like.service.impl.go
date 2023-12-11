package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type LikeServiceImpl struct {
	likecollection *mongo.Collection
	ctx            context.Context
}

func NewLikeService(likecollection *mongo.Collection, ctx context.Context) LikeService {
	return &LikeServiceImpl{
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
