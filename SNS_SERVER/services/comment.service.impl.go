package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/utils"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/bson"
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

func (c *CommentServiceImpl) CreateComment(request *dto.CreateCommentRequest, sessionInfo *dto.SessionInfo) error {
	ctx := context.Background()
	sessionDataJSON, err := c.getSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		return err
	}
	var sessionData map[string]string
	if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
		return err
	}
	userID := sessionData["user_id"]

	var comment models.Comment
	comment.ID = uuid.NewString()
	comment.AuthorID = userID
	comment.PostID = request.PostID
	comment.Content = request.Content
	comment.CreatedAt = time.Now()
	comment.UpdatedAt = time.Now()
	_, err = c.commentcollection.InsertOne(c.ctx, comment)
	if err != nil {
		return err
	}
	return nil
}

func (c *CommentServiceImpl) GetComments(getCommentRequest *dto.GetCommentRequest) ([]models.Comment, error) {
	filter := bson.M{}

	switch getCommentRequest.Type {
	case 1:
		filter["comment_post_id"] = bson.M{"$regex": getCommentRequest.Body, "$options": "i"}
	case 2:
		filter["comment_author_id"] = bson.M{"$regex": getCommentRequest.Body, "$options": "i"}
	default:
		return nil, errors.New("GET_POSTS_ERROR: invalid search type")
	}

	cursor, err := c.commentcollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var comments []models.Comment
	if err = cursor.All(context.Background(), &comments); err != nil {
		return nil, err
	}
	return comments, nil
}

func (c *CommentServiceImpl) UpdateComment(comment models.Comment) error {
	fmt.Println(comment.ID)
	filter := bson.M{"comment_id": comment.ID}

	// 업데이트할 내용 정의
	update := bson.M{
		"$set": bson.M{
			"comment_content":    comment.Content,
			"comment_updated_at": time.Now(),
		},
	}

	// MongoDB에 업데이트 요청
	result, err := c.commentcollection.UpdateOne(c.ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount != 1 {
		return errors.New("UPDATE_POST_ERROR: no matched document found for update")
	}

	return nil
}

func (c *CommentServiceImpl) DeleteComment(commentId string) error {
	filter := bson.D{bson.E{Key: "comment_id", Value: commentId}}
	result, _ := c.commentcollection.DeleteOne(c.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("DELETE_USER_ERROR: no matched document found for delete")
	}
	return nil
}

func (c *CommentServiceImpl) DeleteComments(commentId string) error {
	filter := bson.D{bson.E{Key: "comment_post_id", Value: commentId}}
	result, _ := c.commentcollection.DeleteMany(c.ctx, filter)
	if result.DeletedCount != 1 {
		return errors.New("DELETE_USER_ERROR: no matched document found for delete")
	}
	return nil
}

// 중복되는 메서드는 common service로 통합하여 리팩토링 필요
func (c *CommentServiceImpl) getSessionFromRedis(ctx context.Context, sessionInfo *dto.SessionInfo) (string, error) {
	var sessionDataJSON string
	var err error
	switch sessionInfo.UserType {
	case "regular":
		sessionDataJSON, err = c.redisclient.Get(ctx, "regular_user_session:"+sessionInfo.SessionUUID).Result()
	case "google":
		sessionDataJSON, err = c.redisclient.Get(ctx, "google_user_session:"+sessionInfo.SessionUUID).Result()
	}
	if err != nil {
		return "", err
	}

	sessionDataJSON, _, err = utils.SeparateSessionDataAndSignature([]byte(sessionDataJSON))
	if err != nil {
		return "", err
	}

	return sessionDataJSON, nil
}

// 중복되는 메서드는 common service로 통합하여 리팩토링 필요
func (c *CommentServiceImpl) CanEditPost(sessionInfo *dto.SessionInfo, commentID string) (bool, error) {
	ctx := context.Background()
	sessionDataJSON, err := c.getSessionFromRedis(ctx, sessionInfo)
	if err != nil {
		return false, err
	}

	var sessionData map[string]string
	if err := json.Unmarshal([]byte(sessionDataJSON), &sessionData); err != nil {
		return false, err
	}
	userID := sessionData["user_id"]

	var comment models.Comment
	if err := c.commentcollection.FindOne(ctx, bson.M{"comment_id": commentID}).Decode(&comment); err != nil {
		return false, err
	}

	return comment.AuthorID == userID, nil
}
