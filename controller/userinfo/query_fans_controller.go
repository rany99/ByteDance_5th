package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFansResponse struct {
	common.CommonResponse
	*userinfo.FansResponse
}

type ProxyQueryFans struct {
	UserId int64  `form:"user_id" validate:"required,numeric,min=1"`
	Token  string `form:"token"   validate:"required,jwt"`
}

func QueryFansController(ctx *gin.Context) {

	// 接收参数
	var p ProxyQueryFans
	err := ctx.ShouldBindQuery(&p)

	// 参数校验
	if err = common.Validate.Struct(p); err != nil {
		QueryFansFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 解析user_id
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		PostFollowFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 调用service层
	fansResponse, err := userinfo.QueryFans(p.UserId, uid)
	if err != nil {
		QueryFansFailed(ctx, err.Error())
		return
	}

	// 封装返回
	QueryFansSucceed(ctx, fansResponse)
}

func QueryFansFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryFansResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func QueryFansSucceed(ctx *gin.Context, fansResponse *userinfo.FansResponse) {
	ctx.JSON(http.StatusOK, QueryFansResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		FansResponse: fansResponse,
	})
}
