package controllers

import (
	"context"
	"net/http"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
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

// CreateComment godoc
//
//	@Summary		write the comment
//	@Tags			CRUD comment
//	@Description	write the comment to mongodb
//	@Description	id string | post id for comment
//	@Description	content string | comment content
//	@Accept			json
//	@Produce		json
//	@Param			comment	body		dto.CreateCommentRequest	true	"Comment Data"
//	@Success		200		{string}	success
//	@Router			/comment/create [post]
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

// GetComments godoc
//
//	@Summary		get post comments
//	@Tags			CRUD comment
//	@Description	get post comments to mongodb
//	@Description	type = 1, 2 | search comments type 1: post id, type 2: comment author
//	@Description	body sring | Check whether the search body is included in the type search results
//	@Accept			json
//	@Produce		json
//	@Param			id	body		dto.GetCommentRequest	true	"post id"
//	@Success		200	{object}	[]models.Comment
//	@Router			/comment/get/ [post]
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

// UpdateComment godoc
//
//	@Summary		update regular user data
//	@Tags			CRUD comment
//	@Description	comment owner update regular user informations to mongodb
//	@Accept			json
//	@Produce		json
//	@Param			comment	request		body	dto.CreateCommentRequest	true	"Comment update request"
//	@Success		200		{string}	success
//	@Router			/comment/update [patch]
func (cc *CommentController) UpdateComment(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var commentRequest dto.CreateCommentRequest
	if err := ctx.ShouldBindJSON(&commentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if canEdit, err := cc.CommentService.CanEditPost(sessionInfo, commentRequest.PostID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := cc.CommentService.UpdateComment(&commentRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeleteComment godoc
//
//	@Summary		delete the comment
//	@Tags			CRUD comment
//	@Description	comment owner delete the comment from DB
//	@Accept			json
//	@Produce		json
//	@Param			comment_id	path		string	true	"comment_id"
//	@Success		200			{string}	success
//	@Router			/comment/delete/{comment_id} [delete]
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
	commentroute.POST("/get", cc.GetComments)
	commentroute.PATCH("/update", cc.UpdateComment)
	commentroute.DELETE("/delete/:comment_id", cc.DeleteComment)
}
