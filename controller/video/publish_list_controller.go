package video

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/video"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryPublishListResponse struct {
	common.CommonResponse
	*video.PublishListResponse
}

type ProxyQueryPublishList struct {
	UserId int64  `form:"user_id" validate:"required,numeric,min=1"`
	Token  string `form:"token"   validate:"required,jwt"`
}

// QueryPublishListController Controller层
func QueryPublishListController(ctx *gin.Context) {

	// 解析uid
	rawId, _ := ctx.Get("user_id")
	uid, ok := rawId.(int64)
	if !ok {
		QueryPublishListFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 绑定数据
	var p ProxyQueryPublishList
	err := ctx.ShouldBindQuery(&p)

	// 校验数据
	err = common.Validate.Struct(p)
	if err != nil {
		QueryPublishListFailed(ctx, errortype.DataNotMatchErr)
		return
	}

	// 调用service层
	publishListResponse, err := video.QueryPublishListByUid(p.UserId, uid)
	if err != nil {
		QueryPublishListFailed(ctx, err.Error())
		return
	}

	// 封装数据
	QueryPublishListSucceed(ctx, publishListResponse)
}

// QueryPublishListSucceed 获取发布列表成功
func QueryPublishListSucceed(ctx *gin.Context, publishListResponse *video.PublishListResponse) {
	ctx.JSON(http.StatusOK, QueryPublishListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		PublishListResponse: publishListResponse,
	})
}

// QueryPublishListFailed 获取发布列表失败
func QueryPublishListFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryPublishListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		PublishListResponse: nil,
	})
}
