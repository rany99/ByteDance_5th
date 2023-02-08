package routers

import (
	"ByteDance_5th/handlers/login"
	"ByteDance_5th/handlers/userinfo"
	"ByteDance_5th/middle"
	"ByteDance_5th/models"
	"github.com/gin-gonic/gin"
)

func RoutersInit() *gin.Engine {
	models.InitDB()
	r := gin.Default()
	r.Static("static", "./static")

	BG := r.Group("/douyin")
	BG.GET("/user/", middle.Permission(), userinfo.UserInfoHandler)
	BG.POST("/user/login/", middle.ShaMiddleWare(), login.UserLoginHandler)
	BG.POST("/user/register/", middle.ShaMiddleWare(), login.RegisterHandler)

	return r
}
