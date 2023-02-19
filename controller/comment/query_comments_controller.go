package comment

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/comment"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type QueryCommentListResponse struct {
	common.CommonResponse
	*comment.QueryCommentListResponse
}

type ProxyQueryCommentList struct {
	VideoId int64 `form:"video_id" validate:"required,numeric,min=1"`
}

func QueryCommentListController(ctx *gin.Context) {

	var p ProxyQueryCommentList
	err := ctx.ShouldBindQuery(&p)

	err = common.Validate.Struct(p)
	log.Println("p", p.VideoId)
	if err != nil {
		QueryCommentListFailed(ctx, errortype.ParseVideoIdErr)
		return
	}

	commentList, err := comment.QueryCommentList(p.VideoId)
	if err != nil {
		QueryCommentListFailed(ctx, err.Error())
	}

	QueryCommentListSucceed(ctx, commentList)
}

// QueryCommentListSucceed 成功
func QueryCommentListSucceed(ctx *gin.Context, commentList *comment.QueryCommentListResponse) {
	ctx.JSON(http.StatusOK, QueryCommentListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		QueryCommentListResponse: commentList,
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
