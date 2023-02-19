package middleware

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util"
	"crypto/sha1"
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// JwtMiddleware 中间件
func JwtMiddleware() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenStr := context.Query("token")
		if tokenStr == "" {
			tokenStr = context.PostForm("token")
		}
		//缺少用户信息
		if tokenStr == "" {
			context.JSON(http.StatusOK, common.CommonResponse{
				StatusCode: 401,
				StatusMsg:  errortype.UserNoAuthenticatedErr,
			})
			context.Abort()
			return
		}
		//验证token
		tokenDecoded, ok := util.DecodeToken(tokenStr)
		if !ok {
			context.JSON(http.StatusOK, common.CommonResponse{
				StatusCode: 403,
				StatusMsg:  errortype.ParseTokenErr,
			})
			context.Abort()
			return
		}
		if time.Now().Unix() > tokenDecoded.ExpiresAt {
			context.JSON(http.StatusOK, common.CommonResponse{
				StatusCode: 402,
				StatusMsg:  errortype.TokenOutDateErr,
			})
			context.Abort()
			return
		}
		context.Set("user_id", tokenDecoded.UserId)
		context.Next()
	}
}

func ShaMiddleWare() gin.HandlerFunc {
	return func(context *gin.Context) {
		password := context.Query("password")
		if password == "" {
			password = context.PostForm("password")
		}
		context.Set("password", Sha1(password))
		context.Next()
	}
}

func NoAuthToGetUserId() gin.HandlerFunc {
	return func(context *gin.Context) {
		rawId := context.Query("user_id")
		if rawId == "" {
			rawId = context.PostForm("user_id")
		}
		//用户未认证
		if rawId == "" {
			context.JSON(http.StatusOK, common.CommonResponse{StatusCode: 401, StatusMsg: errortype.UserNoExistErr})
			context.Abort() //阻止执行
			return
		}
		userid, err := strconv.ParseInt(rawId, 10, 64)
		if err != nil {
			context.JSON(http.StatusOK, common.CommonResponse{StatusCode: 401, StatusMsg: errortype.UserNoExistErr})
			context.Abort()
		}
		context.Set("user_id", userid)
		context.Next()
	}
}

func Sha1(str string) string {
	o := sha1.New()
	o.Write([]byte(str))
	return hex.EncodeToString(o.Sum(nil))
}
