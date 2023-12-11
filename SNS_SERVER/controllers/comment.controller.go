package controllers

import (
	"context"
	"net/http"

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
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) GetComments(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) UpdateComment(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) DeleteComment(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (cc *CommentController) RegisterCommentRoutes(rg *gin.RouterGroup) {
	commentroute := rg.Group("/comment")
	commentroute.POST("/create", cc.CreateComment)
	commentroute.GET("/get", cc.GetComments)
	commentroute.PATCH("/update", cc.UpdateComment)
	commentroute.DELETE("/delete", cc.DeleteComment)
}
