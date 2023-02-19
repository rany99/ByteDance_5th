package message

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/server/message"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type PostMessageResponse struct {
	common.CommonResponse
}

type ProxyPostMessageController struct {
	fromId     int64
	toId       int64
	actionType int64
	content    string

	*gin.Context
}

func PostMessageController(ctx *gin.Context) {
	NewProxyPostMessageController(ctx).Operation()
}

func NewProxyPostMessageController(ctx *gin.Context) *ProxyPostMessageController {
	return &ProxyPostMessageController{Context: ctx}
}

func (p *ProxyPostMessageController) Operation() {
	if err := p.ParseJSON(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	if err := message.PostMessage(p.fromId, p.toId, p.actionType, p.content); err != nil {
		p.SendFailed(err.Error())
	}
	p.SendSucceed()
}

func (p *ProxyPostMessageController) ParseJSON() error {
	//解析user_id
	rawFromId, _ := p.Get("user_id")
	fromId, ok := rawFromId.(int64)
	if !ok {
		errors.New(errortype.ParseMsgFromUserIdErr)
	}
	//解析to_user_id
	rawToId := p.Query("to_user_id")
	toId, err := strconv.ParseInt(rawToId, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseMsgToUserIdErr)
	}
	//解析action_type
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseActionTypeErr)
	}
	//解析content
	content := p.Query("content")
	//填入代理层
	p.fromId, p.toId, p.actionType, p.content = fromId, toId, actionType, content
	return nil
}

func (p *ProxyPostMessageController) SendFailed(msg string) {
	p.JSON(http.StatusOK, PostMessageResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *ProxyPostMessageController) SendSucceed() {
	p.JSON(http.StatusOK, PostMessageResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "发送成功",
		}})
}
