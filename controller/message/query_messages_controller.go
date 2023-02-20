package message

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/message"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

type QueryMessageListResponse struct {
	common.CommonResponse
	*message.MessagesResponse
}

type ProxyQueryMessageList struct {
	Token      string `form:"token"         validate:"required,jwt"`
	ToUserId   int64  `form:"to_user_id"    validate:"required,numeric,min=1"`
	PreMsgTime int64  `form:"pre_msg_time"`
}

func QueryMessageListController(ctx *gin.Context) {

	// 绑定参数
	var p ProxyQueryMessageList
	err := ctx.ShouldBindQuery(&p)

	// 校验参数
	err = common.Validate.Struct(p)
	if err != nil {
		QueryMessageListFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 解析from_user_id
	rawFromUid, _ := ctx.Get("user_id")
	fromId, ok := rawFromUid.(int64)
	if !ok {
		QueryMessageListFailed(ctx, errortype.ParseMsgFromUserIdErr)
		return
	}

	if p.PreMsgTime == 0 {
		p.PreMsgTime = time.Now().Unix()
	}
	log.Println("PreMsgTime:", p.PreMsgTime)

	// 调用service层
	messagesResponse, err := message.QueryMessageList(fromId, p.ToUserId, p.PreMsgTime)
	if err != nil {
		QueryMessageListFailed(ctx, err.Error())
		return
	}

	QueryMessageListSucceed(ctx, messagesResponse)
}

// QueryMessageListFailed 查询失败
func QueryMessageListFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryMessageListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// QueryMessageListSucceed 查询成功
func QueryMessageListSucceed(ctx *gin.Context, messagesResponse *message.MessagesResponse) {
	ctx.JSON(http.StatusOK, QueryMessageListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		MessagesResponse: messagesResponse,
	})
}
