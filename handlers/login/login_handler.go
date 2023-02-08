package login

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/login"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserLoginResponse struct {
	models.CommonResponse
	*login.LoginResponse
}

func UserLoginHandler(c *gin.Context) {
	username := c.Query("username")
	raw, _ := c.Get("password")
	password, ok := raw.(string)
	if !ok {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{
				StatusCode: 1,
				StatusMsg:  "密码解析错误",
			},
		})
	}
	//log.Println("密码解析成功")
	userLoginResponse, err := login.QueryUserLogin(username, password)

	//用户不存在返回对应的错误
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: err.Error()},
		})
		return
	}
	//log.Println("userLoginResponse.Token:", userLoginResponse.Token)
	//log.Println("userLoginResponse.UserId:", userLoginResponse.UserId)
	//用户存在，返回相应的id和token
	c.JSON(http.StatusOK, UserLoginResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "登陆成功",
		},
		LoginResponse: userLoginResponse,
	})
}
