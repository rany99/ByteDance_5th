package user

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/user"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterController(ctx *gin.Context) {
	//username := ctx.Query("username")
	//raw, _ := ctx.Get("password")
	//password, ok := raw.(string)

	var p ProxyUser

	err := ctx.ShouldBindQuery(&p)
	err = common.Validate.Struct((p))
	if err != nil {
		RegisterFailed(ctx, errortype.ParsePasswordErr)
		return
	}

	registerResponse, err := user.PostUserLogin(p.Username, p.Password)
	if err != nil {
		RegisterFailed(ctx, err.Error())
		return
	}

	RegisterSucceed(ctx, registerResponse)
}

func RegisterFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func RegisterSucceed(ctx *gin.Context, registerResponse *user.LoginResponse) {
	ctx.JSON(http.StatusOK, LoginResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		LoginResponse: registerResponse,
	})
}
