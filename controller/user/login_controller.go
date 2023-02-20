package user

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginResponse struct {
	common.CommonResponse
	*user.LoginResponse
}

func LoginController(ctx *gin.Context) {

	username := ctx.Query("username")
	rawPassword, _ := ctx.Get("password")
	password, ok := rawPassword.(string)
	if !ok {
		LoginFailed(ctx, errortype.ParsePasswordErr)
		return
	}

	userLoginResponse, err := user.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		LoginFailed(ctx, err.Error())
		return
	}
	//log.Println("userLoginResponse.Token:", userLoginResponse.Token)
	//log.Println("userLoginResponse.UserId:", userLoginResponse.UserId)
	//用户存在，返回相应的id和token

	LoginSucceed(ctx, userLoginResponse)
}

func LoginFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func LoginSucceed(ctx *gin.Context, userLoginResponse *user.LoginResponse) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "登录成功",
		},
		LoginResponse: userLoginResponse,
	})
}
