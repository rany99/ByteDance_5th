package userinfo

import (
	"ByteDance_5th/cmd/models"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryUserInfoResponse struct {
	common.CommonResponse
	*models.UserInfo `json:"user"`
}

func QueryUserInfoController(ctx *gin.Context) {

	// 解析uid
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		QueryUserInfoFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 调用models层
	userinfoDAO := models.NewUserInfoDAO()
	var userInfo models.UserInfo
	err := userinfoDAO.QueryUserInfoById(uid, &userInfo)
	if err != nil {
		QueryUserInfoFailed(ctx, err.Error())
		return
	}

	QueryUserInfoSucceed(ctx, &userInfo)
}

// QueryUserInfoFailed 查询失败
func QueryUserInfoFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryUserInfoResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg},
	})
}

// QueryUserInfoSucceed 查询成功
func QueryUserInfoSucceed(ctx *gin.Context, userInfo *models.UserInfo) {
	ctx.JSON(http.StatusOK, QueryUserInfoResponse{
		CommonResponse: common.CommonResponse{StatusCode: 0},
		UserInfo:       userInfo,
	})
}
