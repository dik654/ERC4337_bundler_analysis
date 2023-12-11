package services

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type PostServiceImpl struct {
	postcollection *mongo.Collection
	ctx            context.Context
}

func NewPostService(postcollection *mongo.Collection, ctx context.Context) PostService {
	return &PostServiceImpl{
		postcollection: postcollection,
		ctx:            ctx,
	}
}

func (p *PostServiceImpl) CreatePost() error {
	return nil
}

func (p *PostServiceImpl) GetAllPosts() error {
	return nil
}

func (p *PostServiceImpl) GetPosts() error {
	return nil
}

func (p *PostServiceImpl) GetUserPosts() error {
	return nil
}

func (p *PostServiceImpl) UpdatePost() error {
	return nil
}

func (p *PostServiceImpl) DeletePost() error {
	return nil
}
