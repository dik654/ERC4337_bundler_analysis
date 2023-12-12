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
	var post models.Post
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.CreatePost(&post); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) GetAllPosts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) GetPosts(ctx *gin.Context) {
	var getPostRequest dto.GetPostRequest
	if err := ctx.ShouldBindJSON(&getPostRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	posts, err := pc.PostService.GetPosts(&getPostRequest)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, posts)
}

func (pc *PostController) GetUserPosts(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) UpdatePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) DeletePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) RegisterPostRoutes(rg *gin.RouterGroup) {
	postroute := rg.Group("/post")
	postroute.POST("/create", pc.CreatePost)
	postroute.GET("/getall", pc.GetAllPosts)
	postroute.GET("/get", pc.GetPosts)
	postroute.GET("/get/:id", pc.GetUserPosts)
	postroute.PATCH("/update", pc.UpdatePost)
	postroute.DELETE("/delete", pc.DeletePost)
}
