package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type FollowsResponse struct {
	models.CommonResponse
	*userinfo.Follows
}

type ProxyQueryFollows struct {
	uid int64
	*userinfo.Follows
	*gin.Context
}

func QueryFollowsController(ctx *gin.Context) {
	NewProxyQueryFollows(ctx).Do()
}

func NewProxyQueryFollows(ctx *gin.Context) *ProxyQueryFollows {
	return &ProxyQueryFollows{
		Context: ctx,
	}
}

func (p *ProxyQueryFollows) Do() {
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
	p.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (p *ProxyQueryFollows) SendSucceed(msg string) {
	p.JSON(http.StatusOK, FollowsResponse{
		CommonResponse: models.CommonResponse{
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
		return errors.New("uid解析错误")
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
