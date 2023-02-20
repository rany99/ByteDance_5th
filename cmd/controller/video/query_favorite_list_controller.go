package video

import (
	"ByteDance_5th/cmd/service/video"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFavoriteListResponse struct {
	common.CommonResponse
	*video.FavoriteListResponse
}

func QueryFavoriteListController(ctx *gin.Context) {

	// 解析uid
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		QueryFavoriteListFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 调用service层
	favoriteListResponse, err := video.QueryFavoriteList(uid)
	if err != nil {
		QueryFavoriteListFailed(ctx, err.Error())
		return
	}

	QueryFavoriteListSucceed(ctx, favoriteListResponse)
}

// QueryFavoriteListSucceed 查询成功
func QueryFavoriteListSucceed(ctx *gin.Context, favoriteListResponse *video.FavoriteListResponse) {
	ctx.JSON(http.StatusOK, QueryFavoriteListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		FavoriteListResponse: favoriteListResponse,
	})
}

// QueryFavoriteListFailed 查询失败
func QueryFavoriteListFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryFavoriteListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		FavoriteListResponse: nil,
	})
}
