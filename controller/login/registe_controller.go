package login

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterResponse struct {
	common.CommonResponse
	*login.LoginResponse
}

func RegisterController(ctx *gin.Context) {
	username := ctx.Query("username")
	raw, _ := ctx.Get("password")
	password, ok := raw.(string)
	if !ok {
		ctx.JSON(http.StatusOK, RegisterResponse{
			CommonResponse: common.CommonResponse{
				StatusCode: 1,
				StatusMsg:  errortype.ParsePasswordErr,
			},
		})
		return
	}
	registerResponse, err := login.PostUserLogin(username, password)
	if err != nil {
		ctx.JSON(http.StatusOK, RegisterResponse{
			CommonResponse: common.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, RegisterResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		LoginResponse: registerResponse,
	})
}
