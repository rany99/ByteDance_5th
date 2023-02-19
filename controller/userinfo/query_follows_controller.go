package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FollowsResponse struct {
	common.CommonResponse
	*userinfo.Follows
}

type ProxyQueryFollows struct {
	uid int64
	*userinfo.Follows
	*gin.Context
}

func QueryFollowsController(ctx *gin.Context) {
	NewProxyQueryFollows(ctx).Operation()
}

func NewProxyQueryFollows(ctx *gin.Context) *ProxyQueryFollows {
	return &ProxyQueryFollows{
		Context: ctx,
	}
}

func (p *ProxyQueryFollows) Operation() {
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

func (p *ProxyQueryFollows) SendFailed(msg string) {
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (p *ProxyQueryFollows) SendSucceed(msg string) {
	p.JSON(http.StatusOK, FollowsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		Follows: p.Follows,
	})
}

func (p *ProxyQueryFollows) ParseJSON() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New(errortype.ParseUserIdErr)
	}
	p.uid = uid
	return nil
}

func (p *ProxyQueryFollows) GetData() error {
	list, err := userinfo.QueryFollowList(p.uid)
	if err != nil {
		return err
	}
	p.Follows = list
	return nil
}
