package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserResponse struct {
	common.CommonResponse
	UserInfo *models.UserInfo `json:"user"`
}

type ProxyUserInfo struct {
	c *gin.Context
}

func NewProxyUserInfo(c *gin.Context) *ProxyUserInfo {
	return &ProxyUserInfo{c: c}
}

func InfoController(ctx *gin.Context) {
	p := NewProxyUserInfo(ctx)
	raw, ok := ctx.Get("user_id")
	if !ok {
		p.SendFailed(errortype.ParseUserIdErr)
	}
	err := p.DoQueryUserInfoByUserId(raw)
	if err != nil {
		p.SendFailed(err.Error())
	}
}

func (p *ProxyUserInfo) DoQueryUserInfoByUserId(rawId interface{}) error {
	userId, ok := rawId.(int64)
	if !ok {
		return errors.New(errortype.ParseUserIdErr)
	}

	userinfoDAO := models.NewUserInfoDAO()

	var userInfo models.UserInfo
	err := userinfoDAO.QueryUserInfoById(userId, &userInfo)
	if err != nil {
		return err
	}
	p.SendSucceed(&userInfo)
	return nil
}

func (p *ProxyUserInfo) SendFailed(msg string) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: common.CommonResponse{StatusCode: 1, StatusMsg: msg},
	})
}

func (p *ProxyUserInfo) SendSucceed(user *models.UserInfo) {
	p.c.JSON(http.StatusOK, UserResponse{
		CommonResponse: common.CommonResponse{StatusCode: 0},
		UserInfo:       user,
	})
}
