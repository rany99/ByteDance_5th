package userinfo

import (
	"ByteDance_5th/cmd/service/userinfo"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryFriendsResponse struct {
	common.CommonResponse
	*userinfo.FriendsResponse
}

type ProxyQueryFriends struct {
	UserId int64 `form:"user_id" validate:"required,numeric,min=1"`
	//Token  string `form:"token"   validate:"required,jwt"`
}

func QueryFriendsController(ctx *gin.Context) {

	// 绑定参数
	var p ProxyQueryFollows
	err := ctx.ShouldBindQuery(&p)

	// 解析uid
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		QueryFriendsFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 调用service层
	friendsResponse, err := userinfo.QueryFriends(uid)
	if err != nil {
		QueryFollowsFailed(ctx, err.Error())
		return
	}

	// 封装数据
	QueryFriendsSucceed(ctx, friendsResponse)
}

// QueryFriendsFailed 查询失败
func QueryFriendsFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryFriendsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// QueryFriendsSucceed 查询成功
func QueryFriendsSucceed(ctx *gin.Context, friendResponse *userinfo.FriendsResponse) {
	ctx.JSON(http.StatusOK, QueryFriendsResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
		},
		FriendsResponse: friendResponse,
	})
}
