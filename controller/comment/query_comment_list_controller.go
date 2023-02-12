package comment

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/comment"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CtListResponse struct {
	models.CommonResponse
	*comment.CList
}

type ProxyCommentListController struct {
	uid int64
	vid int64

	*gin.Context
}

func QueryCommentListController(ctx *gin.Context) {
	NewProxyCommentListHandler(ctx).Do()
}

func (p *ProxyCommentListController) Do() {
	if err := p.ParseJson(); err != nil {
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

func (p *ProxyCommentListController) ParseJson() error {
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New("uid解析错误")
	}
	vid, err := strconv.ParseInt(p.Query("video_id"), 10, 64)
	if err != nil {
		return err
	}
	p.uid, p.vid = uid, vid
	return nil
}

func (p *ProxyCommentListController) SendFailed(msg string) {
	p.JSON(http.StatusOK, CtListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}

func (p *ProxyCommentListController) SendSucceed(commentList *comment.CList) {
	p.JSON(http.StatusOK, CtListResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		CList: commentList,
	})
}

func NewProxyCommentListHandler(ctx *gin.Context) *ProxyCommentListController {
	return &ProxyCommentListController{
		Context: ctx,
	}
}
