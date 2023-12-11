package controllers

import (
	"context"
	"net/http"

	"github.com/dik654/Go_projects/SNS_SERVER/services"
	"github.com/gin-gonic/gin"
)

type LikeController struct {
	LikeService services.LikeService
	ctx         context.Context
}

func NewLikeController(likeservice services.LikeService) LikeController {
	return LikeController{
		LikeService: likeservice,
	}
}

func (lc *LikeController) LikePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (lc *LikeController) UnLikePost(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (lc *LikeController) GetLikes(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (lc *LikeController) RegisterLikeRoutes(rg *gin.RouterGroup) {
	likeroute := rg.Group("/like")
	likeroute.POST("/up", lc.LikePost)
	likeroute.POST("/down", lc.UnLikePost)
	likeroute.GET("/get", lc.GetLikes)
}
