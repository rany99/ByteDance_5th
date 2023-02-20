package video

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProxyPostFavor struct {
	Token      string `form:"token"        validate:"required,jwt"`
	VideoId    int64  `form:"video_id"     validate:"required,numeric,min=1"`
	ActionType int64  `form:"action_type"  validate:"required,numeric,oneof=1 2"`
}

// PostFavorController 点赞操作
func PostFavorController(ctx *gin.Context) {

	// 获取参数
	var p ProxyPostFavor
	err := ctx.ShouldBindQuery(&p)

	// 校验参数
	if err = common.Validate.Struct(p); err != nil {
		PostFavorFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 解析user_id
	RawUserId, _ := ctx.Get("user_id")
	userId, ok := RawUserId.(int64)
	if ok != true {
		PostFavorFailed(ctx, errortype.ParseUserIdErr)
		return
	}
	//log.Println("userId解析成功：", userId)

	// 调用service层
	if err = video.PostFavor(userId, p.VideoId, p.ActionType); err != nil {
		PostFavorFailed(ctx, err.Error())
		return
	}

	// 点赞成功
	PostFavorSucceed(ctx)
}

// PostFavorFailed 生成错误返回
func PostFavorFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// PostFavorSucceed 生成正确返回
func PostFavorSucceed(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "",
	})
}
