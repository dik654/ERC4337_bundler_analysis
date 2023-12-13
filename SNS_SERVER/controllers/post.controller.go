package controllers

import (
	"context"
	"net/http"

	"github.com/dik654/Go_projects/SNS_SERVER/controllers/dto"
	"github.com/dik654/Go_projects/SNS_SERVER/models"
	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
)

type PostController struct {
	PostService services.PostService
	ctx         context.Context
}

func NewPostController(postservice services.PostService) PostController {
	return PostController{
		PostService: postservice,
	}
}

func (pc *PostController) CreatePost(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var post dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.CreatePost(&post, sessionInfo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	var paginationRequest dto.PaginationRequest
	if err := ctx.ShouldBindJSON(&paginationRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	posts, err := pc.PostService.GetAllPosts(&paginationRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func (pc *PostController) GetPosts(ctx *gin.Context) {
	var getPostRequest dto.GetPostRequest
	if err := ctx.ShouldBindJSON(&getPostRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	posts, err := pc.PostService.GetPosts(&getPostRequest)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	postID := ctx.Param("post_id")
	if canEdit, err := pc.PostService.CanEditPost(sessionInfo, postID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	post.ID = postID
	if err := pc.PostService.UpdatePost(post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	postID := ctx.Param("post_id")
	if canEdit, err := pc.PostService.CanEditPost(sessionInfo, postID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.DeletePost(postID); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) RegisterPostRoutes(rg *gin.RouterGroup) {
	postroute := rg.Group("/post")
	postroute.POST("/create", pc.CreatePost)
	postroute.GET("/getall", pc.GetAllPosts)
	postroute.GET("/get", pc.GetPosts)
	postroute.PATCH("/update/:post_id", pc.UpdatePost)
	postroute.DELETE("/delete/:post_id", pc.DeletePost)
}
