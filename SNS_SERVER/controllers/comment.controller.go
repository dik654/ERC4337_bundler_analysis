package controllers

import (
	"context"
	"net/http"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
)

type CommentController struct {
	CommentService services.CommentService
	ctx            context.Context
}

func NewCommentController(commentservice services.CommentService) CommentController {
	return CommentController{
		CommentService: commentservice,
	}
}

func (cc *CommentController) CreateComment(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var comment dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := cc.CommentService.CreateComment(&comment, sessionInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) GetComments(ctx *gin.Context) {
	var getCommentRequest dto.GetCommentRequest
	if err := ctx.ShouldBindJSON(&getCommentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	comments, err := cc.CommentService.GetComments(&getCommentRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}

func (cc *CommentController) UpdateComment(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	commentID := ctx.Param("comment_id")
	if canEdit, err := cc.CommentService.CanEditPost(sessionInfo, commentID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var comment models.Comment
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	comment.ID = commentID
	if err := cc.CommentService.UpdateComment(comment); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) DeleteComment(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	commentID := ctx.Param("comment_id")
	if canEdit, err := cc.CommentService.CanEditPost(sessionInfo, commentID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := cc.CommentService.DeleteComment(commentID); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) RegisterCommentRoutes(rg *gin.RouterGroup) {
	commentroute := rg.Group("/comment")
	commentroute.POST("/create", cc.CreateComment)
	commentroute.GET("/get", cc.GetComments)
	commentroute.PATCH("/update/:comment_id", cc.UpdateComment)
	commentroute.DELETE("/delete/:comment_id", cc.DeleteComment)
}
