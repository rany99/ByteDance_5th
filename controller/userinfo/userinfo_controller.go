package userinfo

import (
	"ByteDance_5th/models"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserResponse struct {
	models.CommonResponse
	UserInfo *models.UserInfo `json:"user"`
}

type ProxyUserInfo struct {
	c *gin.Context
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

func UserInfoHandler(ctx *gin.Context) {
	p := NewProxyUserInfo(ctx)
	raw, ok := ctx.Get("user_id")
	if !ok {
		p.UserInfoError("userId解析错误")
	}
	err := p.DoQueryUserInfoByUserId(raw)
	if err != nil {
		p.UserInfoError(err.Error())
	}
}

func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New("解析userId失败")
	}
	//由于得到userinfo不需要组装model层的数据，所以直接调用model层的接口
	userinfoDAO := models.NewUserInfoDAO()

	var userInfo models.UserInfo
	err := userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.UserInfoOk(&userInfo)
	return nil
}

func (p *ProxyUserInfo) UserInfoError(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: models.CommonResponse{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) UserInfoOk(user *models.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		UserInfo:       user,
	})
}
