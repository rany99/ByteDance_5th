package routers

import (
	"ByteDance_5th/controller/login"
	"ByteDance_5th/controller/userinfo"
	"ByteDance_5th/controller/video"
	"ByteDance_5th/middle"
	"ByteDance_5th/models"
	"github.com/gin-gonic/gin"
)

func RoutersInit() *gin.Engine {
	models.InitDB()
	r := gin.Default()
	r.Static("static", "./public")

	BG := r.Group("/douyin")
	BG.GET("/feed/", video.FeedListHandler)
	BG.GET("/user/", middle.Permission(), userinfo.UserInfoHandler)
	BG.POST("/user/login/", middle.ShaMiddleWare(), login.UserLoginHandler)
	BG.POST("/user/register/", middle.ShaMiddleWare(), login.RegisterHandler)
	BG.POST("/publish/action/", middle.Permission(), video.PublishHandler)
	BG.GET("/publish/list/", middle.NoAuthToGetUserId(), video.QueryVideoListController)
	return r
}
