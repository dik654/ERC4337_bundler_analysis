package services

import (
	"context"
	"errors"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"go.mongodb.org/mongo-driver/bson"
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

func (p *PostServiceImpl) CreatePost(post *models.Post) error {
	_, err := p.postcollection.InsertOne(p.ctx, post)
	return err
}

func (p *PostServiceImpl) GetAllPosts() error {
	return nil
}

func (p *PostServiceImpl) GetPosts(getPostRequest *dto.GetPostRequest) ([]models.Post, error) {
	filter := bson.M{}

	switch getPostRequest.Type {
	case 1: // 제목으로 검색
		filter["title"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	case 2: // 내용으로 검색
		filter["content"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	case 3: // 작성자로 검색
		filter["author"] = bson.M{"$regex": getPostRequest.Body, "$options": "i"}
	default:
		return nil, errors.New("GET_POSTS_ERROR: invalid search type")
	}

	cursor, err := p.postcollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var posts []models.Post
	if err = cursor.All(context.Background(), &posts); err != nil {
		return nil, err
	}
	return posts, nil
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
