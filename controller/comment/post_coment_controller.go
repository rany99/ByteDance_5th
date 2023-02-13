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

//发布评论
//controller层
//解析JSON信息，并将信息送入server层

type PostCommentResponse struct {
	models.CommonResponse
	*comment.CResponse
}

type ProxyPostCommentController struct {
	uid         int64  //用户ID
	vid         int64  //视频ID
	commentID   int64  //评论ID
	actionType  int64  //操作类型
	commentText string //评论内容

	*gin.Context
}

func PostCommentController(ctx *gin.Context) {
	NewProxyPostCommentController(ctx).Operation()
}

func NewProxyPostCommentController(ctx *gin.Context) *ProxyPostCommentController {
	return &ProxyPostCommentController{Context: ctx}
}

func (p *ProxyPostCommentController) Operation() {
	//解析
	if err := p.ParseJson(); err != nil {
		p.SendFailed(err.Error())
		return
	}

	//发布评论
	ret, err := comment.PostComment(p.uid, p.vid, p.commentID, p.actionType, p.commentText)
	//if ret == nil {
	//	log.Println("PostComment返回的ret为空")
	//}
	if err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSucceed(ret)
}

// ParseJson 对JSON进行解析
func (p *ProxyPostCommentController) ParseJson() error {
	//解析UID
	rawUid := p.Query("user_id")
	uid, err := strconv.ParseInt(rawUid, 10, 64)
	//log.Println("parseJson:uid", uid)
	if err != nil {
		return errors.New("uid解析错误")
	}
	//log.Println("parseJson:uid", uid)

	//解析VID
	rawVid := p.Query("video_id")
	vid, err := strconv.ParseInt(rawVid, 10, 64)
	if err != nil {
		return errors.New("vid解析错误")
	}
	//log.Println("parseJson:vid", vid)

	//解析actionType
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return err
	}
	log.Println("parseJson:actionType", actionType)

	//根据actionType进行相应操作
	switch actionType {
	case 1: //添加评论
		p.commentText = p.Query("comment_text")
		log.Println(p.commentText)
	case 2: //删除评论
		p.commentID, err = strconv.ParseInt(p.Query("comment_id"), 10, 64)
		if err != nil {
			return err
		}
	default:
		return errors.New("action_type只能为1或2")
	}

	//填入防火层
	p.uid, p.vid, p.actionType = uid, vid, actionType
	return nil
}

// SendFailed 失败
func (p *ProxyPostCommentController) SendFailed(msg string) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
		CResponse: nil,
	})
}

// SendSucceed 成功
func (p *ProxyPostCommentController) SendSucceed(comment *comment.CResponse) {
	p.JSON(http.StatusOK, PostCommentResponse{
		CommonResponse: models.CommonResponse{StatusCode: 0},
		CResponse:      comment,
	})
	log.Println("SendSucceed")
}
