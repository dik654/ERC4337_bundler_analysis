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

// CreatePost godoc
//
//	@Summary		upload the post
//	@Tags			CRUD post
//	@Description	write post data to mongodb
//	@Accept			json
//	@Produce		json
//	@Param			post_request	body		dto.CreatePostRequest	true	"Post request"
//	@Success		200				{string}	success
//	@Router			/post/create [post]
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

// GetAllPosts godoc
//
//	@Summary		get all posts
//	@Tags			CRUD post
//	@Description	get all posts by pagination
//	@Description	page_num int64 | page number
//	@Description	page_size int64 | number of post in page
//	@Accept			json
//	@Produce		json
//	@Param			page_request	body		dto.PaginationRequest	true	"Page request"
//	@Success		200				{string}	success
//	@Router			/post/getall [post]
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

// GetPosts godoc
//
//	@Summary		get posts
//	@Tags			CRUD post
//	@Description	get posts from DB
//	@Description	type = 1, 2, 3 | search posts type 1: title, type 2: content, type 3: author
//	@Description	body = string | Check whether the search body is included in the type search results
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.GetPostRequest	true	"Post request"
//	@Success		200		{object}	models.User
//	@Router			/post/get [POST]
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

// UpdateUser godoc
//
//	@Summary		update the post
//	@Tags			CRUD post
//	@Description	post owner update the post
//	@Description	parameter post_id | specific post id
//	@Description	title string | title you want to change
//	@Description	content string | content you want to change
//	@Accept			json
//	@Produce		json
//	@Param			post	body		dto.CreatePostRequest	true	"Update data"
//	@Param			post_id	path		string					true	"post_id"
//	@Success		200		{string}	success
//	@Router			/post/update/{post_id} [patch]
func (pc *PostController) UpdatePost(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	postID := ctx.Param("post_id")
	if canEdit, err := pc.PostService.CanEditPost(ctx, sessionInfo, postID); canEdit != true {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "no authentication to update post"})
		return
	} else if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var post dto.CreatePostRequest
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := pc.PostService.UpdatePost(postID, &post); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// DeletePost godoc
//
//	@Summary		delete the post
//	@Tags			CRUD post
//	@Description	delete the post from DB
//	@Description	parameter post_id | specific post id
//	@Accept			json
//	@Produce		json
//	@Param			post_id	path		string	true	"post_id"
//	@Success		200		{string}	success
//	@Router			/post/delete/{post_id} [delete]
func (pc *PostController) DeletePost(ctx *gin.Context) {
	sessionInfo, err := GetSessionInfo(ctx)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	postID := ctx.Param("post_id")
	if canEdit, err := pc.PostService.CanEditPost(ctx, sessionInfo, postID); canEdit != true || err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.DeletePost(postID); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// LikePost godoc
//
//	@Summary		sign up regular user
//	@Tags			CRUD post
//	@Description	write user informations to mongodb
//	@Description	id string | specific post id
//	@Description	like bool | true: like it, false: unlike it
//	@Accept			json
//	@Produce		json
//	@Param			like_request	body		dto.PostLikeRequest	true	"Like request"
//	@Success		200				{string}	success
//	@Router			/post/like [post]
func (pc *PostController) LikePost(ctx *gin.Context) {
	var postLikeRequest dto.PostLikeRequest
	if err := ctx.ShouldBindJSON(&postLikeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.LikePost(&postLikeRequest); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

// JudgePost godoc
//
//	@Summary		report the post
//	@Tags			CRUD post
//	@Description	write report of the post to DB
//	@Description	type uint8 | type 1: insult, type 2: sexual, type 3: violent
//	@Description	id	string | specific violated post id
//	@Description	body string | details of violation
//	@Accept			json
//	@Produce		json
//	@Param			judge_request	body		models.Judge	true	"Judge request"
//	@Success		200				{string}	success
//	@Router			/post/judge [post]
func (pc *PostController) JudgePost(ctx *gin.Context) {
	var postJudgeRequest models.Judge
	if err := ctx.ShouldBindJSON(&postJudgeRequest); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if err := pc.PostService.JudgePost(postJudgeRequest); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (pc *PostController) RegisterPostRoutes(rg *gin.RouterGroup) {
	postroute := rg.Group("/post")
	postroute.POST("/create", pc.CreatePost)
	postroute.POST("/getall", pc.GetAllPosts)
	postroute.POST("/get", pc.GetPosts)
	postroute.PATCH("/update/:post_id", pc.UpdatePost)
	postroute.DELETE("/delete/:post_id", pc.DeletePost)
	postroute.POST("/like", pc.LikePost)
	postroute.POST("/judge", pc.JudgePost)
}
