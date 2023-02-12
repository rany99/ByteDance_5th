package comment

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/comment"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type PostCommentResponse struct {
	models.CommonResponse
	*comment.CResponse
}

type ProxyPostCommentController struct {
	uid         int64 //用户ID
	vid         int64 //视频ID
	commentID   int64 //评论ID
	actionType  int64 //操作类型
	commentText string

	*gin.Context
}

func PostCommentController(ctx *gin.Context) {
	NewProxyPostCommentController(ctx).Do()
}

func NewProxyPostCommentController(ctx *gin.Context) *ProxyPostCommentController {
	return &ProxyPostCommentController{Context: ctx}
}

func (p *ProxyPostCommentController) Do() {
	if err := p.parseJson(); err != nil {
		p.SendFailed(err.Error())
		return
	}

	ret, err := comment.PostComment(p.uid, p.vid, p.commentID, p.actionType, p.commentText)
	if ret == nil {
		log.Println("PostComment返回的ret为空")
	}
	if err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSucceed(ret)
}

func (p *ProxyPostCommentController) parseJson() error {
	//解析UID
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	log.Println("parseJson:uid", uid)
	if !ok {
		return errors.New("uid解析错误")
	}
	p.uid = uid
	log.Println("parseJson:uid", uid)
	//解析VID
	rawVid := p.Query("video_id")
	vid, err := strconv.ParseInt(rawVid, 10, 64)
	if err != nil {
		return errors.New("vid解析错误")
	}
	p.vid = vid
	log.Println("parseJson:vid", vid)
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return err
	}
	log.Println("parseJson:actionType", actionType)
	switch actionType {
	case 1:
		p.commentText = p.Query("comment_text")
		log.Println(p.commentText)
	case 2:
		p.commentID, err = strconv.ParseInt(p.Query("comment_id"), 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.New("action_type只能为1或2")
	}
	p.actionType = actionType
	return nil
}

func (p *ProxyPostCommentController) SendFailed(msg string) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		CResponse: nil,
	})
}

func (p *ProxyPostCommentController) SendSucceed(comment *comment.CResponse) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		CResponse:      comment,
	})
	log.Println("SendSucceed")
}
