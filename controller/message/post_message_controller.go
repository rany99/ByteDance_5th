package message

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/message"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProxyPostMessage struct {
	Token      string `form:"token"         validate:"required,jwt"`
	ToUserId   int64  `form:"to_user_id"    validate:"required,numeric,min=1"`
	ActionType int64  `form:"action_type"   validate:"required,numeric,oneof=1"`
	Content    string `form:"content"       validate:"required"`
}

func PostMessageController(ctx *gin.Context) {

	// 绑定数据
	var p ProxyPostMessage
	err := ctx.ShouldBindQuery(&p)

	// 参数校验
	if err = common.Validate.Struct(p); err != nil {
		PostMessageFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 解析from_id
	rawFromId, _ := ctx.Get("user_id")
	fromId, ok := rawFromId.(int64)
	if !ok {
		errors.New(errortype.ParseMsgFromUserIdErr)
	}

	// 调用service层
	if err = message.PostMessage(fromId, p.ToUserId, p.ActionType, p.Content); err != nil {
		PostMessageFailed(ctx, err.Error())
		return
	}

	// 调用成功
	PostMessageSucceed(ctx)

}

// PostMessageFailed 发送成功
func PostMessageFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// PostMessageSucceed 发送失败
func PostMessageSucceed(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
	})
}
