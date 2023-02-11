package routers

import (
	"ByteDance_5th/controller/login"
	"ByteDance_5th/controller/userinfo"
	"ByteDance_5th/controller/video"
	"ByteDance_5th/middle"
	"ByteDance_5th/models"
	"github.com/gin-gonic/gin"
)

func DouSheng_RoutersInit() *gin.Engine {

	models.InitDB()

	r := gin.Default()
	r.Static("static", "./public")

	BG := r.Group("/douyin")

	//基础接口
	//视频流接口
	BG.GET("/feed/", video.FeedListHandler)
	//用户注册
	BG.POST("/user/register/", middle.ShaMiddleWare(), login.RegisterHandler)
	//用户登录
	BG.POST("/user/login/", middle.ShaMiddleWare(), login.UserLoginHandler)
	//用户信息
	BG.GET("/user/", middle.Permission(), userinfo.UserInfoHandler)
	//投稿接口
	BG.POST("/publish/action/", middle.Permission(), video.PublishHandler)
	//发布列表
	BG.GET("/publish/list/", middle.NoAuthToGetUserId(), video.QueryVideoListController)

	//互动接口
	//PostFavorHandler 点赞
	BG.POST("/favorite/action/", middle.Permission(), video.PostFavorHandler)
	return r
}
