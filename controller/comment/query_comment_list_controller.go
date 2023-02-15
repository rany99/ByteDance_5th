package comment

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/server/comment"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CtListResponse struct {
	common.CommonResponse
	*comment.CList
}

type ProxyCommentListController struct {
	uid int64
	vid int64

	*gin.Context
}

func QueryCommentListController(ctx *gin.Context) {
	NewProxyCommentListHandler(ctx).Operation()
}

func (p *ProxyCommentListController) Operation() {
	if err := p.ParseJSON(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	cList, err := comment.QueryCommentList(p.uid, p.vid)
	if err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSucceed(cList)
}

// ParseJSON 解析JSON信息
func (p *ProxyCommentListController) ParseJSON() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New(errortype.ParseUserIdErr)
	}
	rawVid := p.Query("video_id")
	vid, err := strconv.ParseInt(rawVid, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseVideoIdErr)
	}

	//填入代理层
	p.uid, p.vid = uid, vid
	return nil
}

// SendFailed 失败
func (p *ProxyCommentListController) SendFailed(msg string) {
	p.JSON(http.StatusOK, CtListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

// SendSucceed 成功
func (p *ProxyCommentListController) SendSucceed(commentList *comment.CList) {
	p.JSON(http.StatusOK, CtListResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CList: commentList,
	})
}

// NewProxyCommentListHandler 创建代理层
func NewProxyCommentListHandler(ctx *gin.Context) *ProxyCommentListController {
	return &ProxyCommentListController{
		Context: ctx,
	}
}
