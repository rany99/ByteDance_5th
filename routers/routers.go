package routers

import (
	"ByteDance_5th/controller/comment"
	"ByteDance_5th/controller/message"
	"ByteDance_5th/controller/user"
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

	BaseGroup := r.Group("/douyin")

	//基础接口
	//视频流接口
	BaseGroup.GET("/feed/", video.QueryFeedListController)
	userGroup := BaseGroup.Group("/user")
	{
		//用户注册
		userGroup.POST("/register/", middleware.ShaMiddleWare(), user.RegisterController)
		//用户登录
		userGroup.POST("/login/", middleware.ShaMiddleWare(), user.LoginController)
		//用户信息
		userGroup.GET("/", middleware.JwtMiddleware(), userinfo.QueryUserInfoController)
	}
	publish := BaseGroup.Group("/publish")
	{
		//投稿接口
		publish.POST("/action/", middleware.JwtMiddleware(), video.PublishController)
		//发布列表
		publish.GET("/list/", middleware.NoAuthToGetUserId(), video.QueryPublishListController)
	}

	//互动接口
	favorite := BaseGroup.Group("/favorite")
	{
		//赞操作
		favorite.POST("/action/", middleware.JwtMiddleware(), video.PostFavorController)
		//喜欢列表
		favorite.GET("/list/", middleware.NoAuthToGetUserId(), video.QueryFavoriteListController)
	}
	commentGroup := BaseGroup.Group("/comment")
	{
		//评论操作
		commentGroup.POST("/action/", middleware.JwtMiddleware(), comment.PostCommentController)
		//评论列表
		commentGroup.GET("/list/", middleware.JwtMiddleware(), comment.QueryCommentListController)
	}

	//社交接口
	relation := BaseGroup.Group("/relation")
	{
		//关注操作
		relation.POST("/action/", middleware.JwtMiddleware(), userinfo.PostFollowController)
		//关注列表
		relation.GET("/follow/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFollowsController)
		//粉丝列表
		relation.GET("/follower/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFansController)
		//朋友列表
		relation.GET("/friend/list/", middleware.NoAuthToGetUserId(), userinfo.QueryFriendsController)
	}

	//消息接口
	messageGroup := BaseGroup.Group("/message")
	{
		//发送消息
		messageGroup.POST("/action/", middleware.JwtMiddleware(), message.PostMessageController)
		//消息记录
		messageGroup.GET("/chat/", middleware.JwtMiddleware(), message.QueryMessageListController)
	}

	return r
}
