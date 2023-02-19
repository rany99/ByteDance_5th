package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFriendsResponse struct {
	common.CommonResponse
	*userinfo.FriendsResponse
}

type ProxyQueryFriendsResponse struct {
	uid int64
	*userinfo.FriendsResponse
	*gin.Context
}

func QueryFriendsController(ctx *gin.Context) {
	NewProxyQueryFriendsResponse(ctx).Operation()
}

func NewProxyQueryFriendsResponse(ctx *gin.Context) *ProxyQueryFriendsResponse {
	return &ProxyQueryFriendsResponse{
		Context: ctx,
	}
}

func (p *ProxyQueryFriendsResponse) Operation() {
	if err := p.ParseJSON(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	if err := p.GetData(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSucceed("查询成功")
}

func (p *ProxyQueryFriendsResponse) ParseJSON() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New(errortype.ParseUserIdErr)
	}
	p.uid = uid
	return nil
}

func (p *ProxyQueryFriendsResponse) GetData() error {
	friends, err := userinfo.QueryFriends(p.uid)
	if err != nil {
		return err
	}
	p.FriendsResponse = friends
	return nil
}

func (p *ProxyQueryFriendsResponse) SendFailed(msg string) {
	p.JSON(http.StatusOK, QueryFriendsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FriendsResponse: nil,
	})
}

func (p *ProxyQueryFriendsResponse) SendSucceed(msg string) {
	p.JSON(http.StatusOK, QueryFriendsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		FriendsResponse: p.FriendsResponse,
	})
}
