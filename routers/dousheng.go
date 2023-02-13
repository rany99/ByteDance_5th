package routers

import (
	"ByteDance_5th/controller/comment"
	"ByteDance_5th/controller/login"
	"ByteDance_5th/controller/userinfo"
	"ByteDance_5th/controller/video"
	"ByteDance_5th/middle"
	"ByteDance_5th/models"
	"github.com/gin-gonic/gin"
)

func DoushengRoutersinit() *gin.Engine {

	//数据库初始化
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
	//赞操作
	BG.POST("/favorite/action/", middle.Permission(), video.PostFavorHandler)
	//喜欢列表
	BG.GET("/favorite/list/", middle.NoAuthToGetUserId(), video.QueryFavoriteListController)
	//评论操作
	BG.POST("/comment/action/", middle.Permission(), comment.PostCommentController)
	//评论列表
	BG.GET("/comment/list/", middle.Permission(), comment.QueryCommentListController)

	//社交接口
	//关注操作
	BG.POST("/relation/action/", middle.Permission(), userinfo.PostFollowController)
	//关注列表
	BG.GET("/relation/follow/list/", middle.NoAuthToGetUserId(), userinfo.QueryFollowsController)
	//粉丝列表
	BG.GET("/relation/follower/list/", middle.NoAuthToGetUserId(), userinfo.QueryFansController)
	return r
}
