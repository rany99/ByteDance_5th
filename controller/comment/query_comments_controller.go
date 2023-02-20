package comment

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/comment"
	"github.com/gin-gonic/gin"
	"net/http"
)

type QueryCommentListResponse struct {
	common.CommonResponse
	*comment.CommentsResponse
}

type ProxyQueryCommentList struct {
	VideoId int64 `form:"video_id" validate:"required,numeric,min=1"`
}

func QueryCommentListController(ctx *gin.Context) {

	// 绑定参数
	var p ProxyQueryCommentList
	err := ctx.ShouldBindQuery(&p)

	// 校验参数
	err = common.Validate.Struct(p)
	if err != nil {
		QueryCommentListFailed(ctx, errortype.ParseVideoIdErr)
		return
	}

	// 调用service层
	commentList, err := comment.QueryCommentList(p.VideoId)
	if err != nil {
		QueryCommentListFailed(ctx, err.Error())
		return
	}

	QueryCommentListSucceed(ctx, commentList)
}

// QueryCommentListSucceed 成功
func QueryCommentListSucceed(ctx *gin.Context, commentList *comment.CommentsResponse) {
	ctx.JSON(http.StatusOK, QueryCommentListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CommentsResponse: commentList,
	})
}

// QueryCommentListFailed 失败
func QueryCommentListFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryCommentListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}
