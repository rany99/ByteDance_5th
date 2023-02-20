package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFollowsResponse struct {
	common.CommonResponse
	*userinfo.FollowsResponse
}

type ProxyQueryFollows struct {
	UserId int64 `form:"user_id" validate:"required,numeric,min=1"`
	//Token  string `form:"token"   validate:"required,jwt"`
}

func QueryFollowsController(ctx *gin.Context) {

	// 接收参数
	var p ProxyQueryFollows
	err := ctx.ShouldBindQuery(&p)

	// 参数校验
	err = common.Validate.Struct(p)
	if err != nil {
		QueryFollowsFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 调用service层
	followsResponse, err := userinfo.QueryFollowList(p.UserId)
	if err != nil {
		QueryFollowsFailed(ctx, err.Error())
		return
	}

	//查询成功
	QueryFollowsSucceed(ctx, followsResponse)
}

func QueryFollowsSucceed(ctx *gin.Context, followsResponse *userinfo.FollowsResponse) {
	ctx.JSON(http.StatusOK, QueryFollowsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		FollowsResponse: followsResponse,
	})
}

func QueryFollowsFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}
