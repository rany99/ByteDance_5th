package user

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterController 用户注册接口
func RegisterController(ctx *gin.Context) {

	username := ctx.Query("username")
	rawPassword, _ := ctx.Get("password")
	password, ok := rawPassword.(string)
	if !ok {
		RegisterFailed(ctx, errortype.ParsePasswordErr)
		return
	}

	// 调用Service层
	registerResponse, err := user.PostUserLogin(username, password)
	if err != nil {
		RegisterFailed(ctx, err.Error())
		return
	}

	// 注册成功
	RegisterSucceed(ctx, registerResponse)
}

// RegisterFailed 注册失败
func RegisterFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// RegisterSucceed 注册成功
func RegisterSucceed(ctx *gin.Context, registerResponse *user.LoginResponse) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		LoginResponse: registerResponse,
	})
}
