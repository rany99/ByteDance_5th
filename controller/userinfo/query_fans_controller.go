package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFansResponse struct {
	common.CommonResponse
	*userinfo.FansResponse
}

type ProxyQueryFansController struct {
	uid int64
	*userinfo.FansResponse
	*gin.Context
}

func QueryFansController(ctx *gin.Context) {
	NewProxyQueryFansController(ctx).Operation()
}

func NewProxyQueryFansController(ctx *gin.Context) *ProxyQueryFansController {
	return &ProxyQueryFansController{Context: ctx}
}

func (p *ProxyQueryFansController) Operation() {
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
		return errors.New(errortype.ParseUserIdErr)
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
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  msg,
		},
		FansResponse: p.FansResponse,
	})
}

func (p *ProxyQueryFansController) SendFailed(msg string) {
	p.JSON(http.StatusOK, QueryFansResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FansResponse: nil,
	})
}
