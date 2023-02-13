package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFansResponse struct {
	models.CommonResponse
	*userinfo.FansResponse
}

type ProxyQueryFansController struct {
	uid int64
	*userinfo.FansResponse
	*gin.Context
}

func QueryFansController(ctx *gin.Context) {
	NewProxyQueryFansController(ctx).Do()
}

func NewProxyQueryFansController(ctx *gin.Context) *ProxyQueryFansController {
	return &ProxyQueryFansController{Context: ctx}
}

func (p *ProxyQueryFansController) Do() {
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

func (p *ProxyQueryFansController) ParseJSON() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New("uid解析错误")
	}
	p.uid = uid
	return nil
}

func (p *ProxyQueryFansController) GetData() error {
	fans, err := userinfo.QueryFans(p.uid)
	if err != nil {
		return err
	}
	p.FansResponse = fans
	return nil
}

func (p *ProxyQueryFansController) SendSucceed(msg string) {
	p.JSON(http.StatusOK, QueryFansResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		FansResponse: p.FansResponse,
	})
}

func (p *ProxyQueryFansController) SendFailed(msg string) {
	p.JSON(http.StatusOK, QueryFansResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FansResponse: nil,
	})
}
