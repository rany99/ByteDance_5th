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

type MsgListResponse struct {
	common.CommonResponse
	*message.MList
}

type ProxyMessageListController struct {
	toId   int64
	fromId int64
	vid    int64

	*gin.Context
}

func QueryMessageListController(ctx *gin.Context) {
	NewProxyMessageListController(ctx).Operation()
}

func NewProxyMessageListController(ctx *gin.Context) *ProxyMessageListController {
	return &ProxyMessageListController{
		Context: ctx,
	}
}
func (p *ProxyMessageListController) Operation() {
	if err := p.ParseJSON(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	//log.Println("f:", p.fromId)
	//log.Println("t", p.toId)
	msgList, err := message.QueryMessageList(p.fromId, p.toId)
	if err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSucceed(msgList)
}

func (p *ProxyMessageListController) ParseJSON() error {
	rawFromUid, _ := p.Get("user_id")
	fromId, ok := rawFromUid.(int64)
	//log.Println("fromId", fromId)
	if !ok {
		return errors.New(errortype.ParseMsgFromUserIdErr)
	}

	rawToUid := p.Query("to_user_id")
	toId, err := strconv.ParseInt(rawToUid, 10, 64)
	//log.Println("to_user_id", toId)
	if err != nil {
		return errors.New(errortype.ParseToUserIdErr)
	}

	p.fromId, p.toId = fromId, toId
	return nil
}

func (p *ProxyMessageListController) SendFailed(msg string) {
	p.JSON(http.StatusOK, MsgListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		MList: nil,
	})
}

func (p *ProxyMessageListController) SendSucceed(msgList *message.MList) {
	p.JSON(http.StatusOK, MsgListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		MList: msgList,
	})
}
