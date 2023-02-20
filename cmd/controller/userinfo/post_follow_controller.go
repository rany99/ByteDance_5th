package userinfo

import (
	"ByteDance_5th/cmd/service/userinfo"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/constantval"
	"ByteDance_5th/pkg/errortype"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProxyPostFollow struct {
	Token      string `form:"token"        validate:"required,jwt"`
	ToUserId   int64  `form:"to_user_id"   validate:"required,numeric,min=1"`
	ActionType int64  `form:"action_type"  validate:"required,numeric,oneof=1 2"`
}

func PostFollowController(ctx *gin.Context) {

	// 接受参数
	var p ProxyPostFollow
	err := ctx.ShouldBindQuery(&p)

	// 解析user_id
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		PostFollowFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 参数校验
	if err = common.Validate.Struct(p); err != nil {
		PostFollowFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 根据 service 返回的状态码返回相应的错误信息
	if stateCode := userinfo.PostFollow(uid, p.ToUserId, p.ActionType); stateCode != constantval.FollowSucceed {
		switch stateCode {
		case constantval.FollowUserNoExist:
			PostFollowFailed(ctx, errortype.FollowUserNoExistErr)
		case constantval.PostFollowActionTypeWrong:
			PostFollowFailed(ctx, errortype.PostFollowActionTypeErr)
		case constantval.CantFollowSelf:
			PostFollowFailed(ctx, errortype.CantFollowSelfErr)
		case constantval.FollowAgain:
			PostFollowFailed(ctx, errortype.FollowAgainErr)
		}
	}

	PostFollowSucceed(ctx)
}

// PostFollowSucceed 关注成功
func PostFollowSucceed(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "关注成功",
	})
}

// PostFollowFailed 关注失败
func PostFollowFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}
