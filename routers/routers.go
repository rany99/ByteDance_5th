package routers

import (
	"ByteDance_5th/controller/comment"
	"ByteDance_5th/controller/login"
	"ByteDance_5th/controller/message"
	"ByteDance_5th/controller/userinfo"
	"ByteDance_5th/controller/video"
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouters() *gin.Engine {

	//数据库初始化
	models.InitDB()

	r := gin.Default()
	r.Static("static", "./public")

	BG := r.Group("/douyin")

	//基础接口
	//视频流接口
	BG.GET("/feed/", video.FeedListController)
	//用户注册
	BG.POST("/user/register/", middleware.ShaMiddleWare(), login.RegisterController)
	//用户登录
	BG.POST("/user/login/", middleware.ShaMiddleWare(), login.UserLoginController)
	//用户信息
	BG.GET("/user/", middleware.Permission(), userinfo.InfoController)
	//投稿接口
	BG.POST("/publish/action/", middleware.Permission(), video.PublishHandler)
	//发布列表
	BG.GET("/publish/list/", middleware.NoAuthToGetUserId(), video.QueryVideoListController)

	//互动接口
	//赞操作
	BG.POST("/favorite/action/", middleware.Permission(), video.PostFavorController)
	//喜欢列表
	BG.GET("/favorite/list/", middleware.NoAuthToGetUserId(), video.QueryFavoriteListController)
	//评论操作
	BG.POST("/comment/action/", middleware.Permission(), comment.PostCommentController)
	//评论列表
	BG.GET("/comment/list/", middleware.Permission(), comment.QueryCommentListController)

	//社交接口
	//关注操作
	BG.POST("/relation/action/", middleware.Permission(), userinfo.PostFollowController)
	//关注列表
	BG.GET("/relation/follow/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFollowsController)
	//粉丝列表
	BG.GET("/relation/follower/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFansController)
	//朋友列表
	BG.GET("/relation/friend/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFriendsController)

	//消息接口
	//消息记录
	BG.GET("/message/chat/", middleware.Permission(), message.QueryMessageListController)

	return r
}
