package comment

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/comment"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PostCommentResponse struct {
	common.CommonResponse
	*comment.CommentResponse
}

// ProxyPostComment 代理层
type ProxyPostComment struct {
	Token       string `form:"token"         validate:"required,jwt"`
	VideoId     int64  `form:"video_id"      validate:"required,numeric,min=1"`
	ActionType  int64  `form:"action_type"   validate:"required,numeric,oneof=1 2"`
	CommentText string `form:"comment_text"`
	CommentId   int64  `form:"comment_id"`
}

func PostCommentController(ctx *gin.Context) {
	var p ProxyPostComment
	// 参数绑定
	err := ctx.ShouldBindQuery(&p)

	// 解析user_id
	rawUid, _ := ctx.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		PostCommentFailed(ctx, errortype.ParseUserIdErr)
		return
	}

	// 参数校验
	if err = common.Validate.Struct(p); err != nil {
		PostCommentFailed(ctx, err.Error())
		return
	}

	// 传入service层
	var postComment *comment.CommentResponse
	postComment, err = comment.PostComment(uid, p.VideoId, p.CommentId, p.ActionType, p.CommentText)
	if err != nil {
		PostCommentFailed(ctx, err.Error())
	}

	// 打包数据
	PostCommentSucceed(ctx, postComment)
}

// PostCommentFailed 失败
func PostCommentFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// PostCommentSucceed 成功
func PostCommentSucceed(ctx *gin.Context, comment *comment.CommentResponse) {
	ctx.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse:  common.CommonResponse{StatusCode: 0},
		CommentResponse: comment,
	})
}
