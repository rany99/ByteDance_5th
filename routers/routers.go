package routers

import (
	comment2 "ByteDance_5th/cmd/controller/comment"
	message2 "ByteDance_5th/cmd/controller/message"
	user2 "ByteDance_5th/cmd/controller/user"
	userinfo2 "ByteDance_5th/cmd/controller/userinfo"
	video2 "ByteDance_5th/cmd/controller/video"
	"ByteDance_5th/cmd/models"
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
	BaseGroup.GET("/feed/", video2.QueryFeedListController)
	userGroup := BaseGroup.Group("/user")
	{
		//用户注册
		userGroup.POST("/register/", middleware.ShaMiddleWare(), user2.RegisterController)
		//用户登录
		userGroup.POST("/login/", middleware.ShaMiddleWare(), user2.LoginController)
		//用户信息
		userGroup.GET("/", middleware.JwtMiddleware(), userinfo2.QueryUserInfoController)
	}
	publish := BaseGroup.Group("/publish")
	{
		//投稿接口
		publish.POST("/action/", middleware.JwtMiddleware(), video2.PublishController)
		//发布列表
		publish.GET("/list/", middleware.NoAuthToGetUserId(), video2.QueryPublishListController)
	}

	//互动接口
	favorite := BaseGroup.Group("/favorite")
	{
		//赞操作
		favorite.POST("/action/", middleware.JwtMiddleware(), video2.PostFavorController)
		//喜欢列表
		favorite.GET("/list/", middleware.NoAuthToGetUserId(), video2.QueryFavoriteListController)
	}
	commentGroup := BaseGroup.Group("/comment")
	{
		//评论操作
		commentGroup.POST("/action/", middleware.JwtMiddleware(), comment2.PostCommentController)
		//评论列表
		commentGroup.GET("/list/", middleware.JwtMiddleware(), comment2.QueryCommentListController)
	}

	//社交接口
	relation := BaseGroup.Group("/relation")
	{
		//关注操作
		relation.POST("/action/", middleware.JwtMiddleware(), userinfo2.PostFollowController)
		//关注列表
		relation.GET("/follow/list/", middleware.NoAuthToGetUserId(), userinfo2.QueryFollowsController)
		//粉丝列表
		relation.GET("/follower/list/", middleware.NoAuthToGetUserId(), userinfo2.QueryFansController)
		//朋友列表
		relation.GET("/friend/list/", middleware.NoAuthToGetUserId(), userinfo2.QueryFriendsController)
	}

	//消息接口
	messageGroup := BaseGroup.Group("/message")
	{
		//发送消息
		messageGroup.POST("/action/", middleware.JwtMiddleware(), message2.PostMessageController)
		//消息记录
		messageGroup.GET("/chat/", middleware.JwtMiddleware(), message2.QueryMessageListController)
	}

	return r
}
