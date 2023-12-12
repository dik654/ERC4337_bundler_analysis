package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type CommentServiceImpl struct {
	commentcollection *mongo.Collection
	ctx               context.Context
}

func NewCommentService(commentcollection *mongo.Collection, ctx context.Context) CommentService {
	return &CommentServiceImpl{
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

func (c *CommentServiceImpl) GetUserComments() error {
	return nil
}

func (c *CommentServiceImpl) UpdateComment() error {
	return nil
}

func (c *CommentServiceImpl) DeleteComment() error {
	return nil
}
