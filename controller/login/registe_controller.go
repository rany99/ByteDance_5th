package login

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterResponse struct {
	models.CommonResponse
	*login.LoginResponse
}

func RegisterHandler(ctx *gin.Context) {
	username := ctx.Query("username")
	raw, _ := ctx.Get("password")
	password, ok := raw.(string)
	if !ok {
		ctx.JSON(http.StatusOK, RegisterResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
		return
	}
	registerResponse, err := login.PostUserLogin(username, password)
	if err != nil {
		ctx.JSON(http.StatusOK, RegisterResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  err.Error(),
			},
		})
		return
	}
	ctx.JSON(http.StatusOK, RegisterResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
		},
		LoginResponse: registerResponse,
	})
}
