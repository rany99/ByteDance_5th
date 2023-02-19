package video

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/video"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func PostFavorController(ctx *gin.Context) {
	NewPostFavorController(ctx).Operation()
}

func NewPostFavorController(ctx *gin.Context) *ProxyPostFavorHandler {
	return &ProxyPostFavorHandler{Context: ctx}
}

type ProxyPostFavorHandler struct {
	*gin.Context
	userId     int64
	videoId    int64
	actionType int64
}

func (p *ProxyPostFavorHandler) Operation() {
	//解码
	if err := p.ParseNum(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	//sever层执行点赞操作
	if err := video.PostFavor(p.userId, p.videoId, p.actionType); err != nil {
		p.SendFailed(err.Error())
		return
	}
	//点赞成功
	p.SendSucceed()
}

// ParseNum 解析UserId
func (p *ProxyPostFavorHandler) ParseNum() error {
	//解析user_id
	RawUserId, _ := p.Get("user_id")
	userId, ok := RawUserId.(int64)
	if ok != true {
		return errors.New(errortype.ParseUserIdErr)
	}
	//log.Println("userId解析成功：", userId)

	//解析video_id
	RawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(RawVideoId, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseVideoIdErr)
	}
	//log.Println("videoId解析成功：", videoId)

	//解析action_type
	RawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(RawActionType, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseActionTypeErr)
	}
	//log.Println("actionType解析成功：", actionType)

	//校验actionType
	if actionType != 1 && actionType != 2 {
		return errors.New(errortype.PostFavorActionTypeErr)
	}

	//填入代理层
	p.videoId, p.actionType, p.userId = videoId, actionType, userId
	return nil
}

// SendFailed 生成错误返回
func (p *ProxyPostFavorHandler) SendFailed(msg string) {
	//log.Println("SendErr")
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// SendSucceed 生成正确返回
func (p *ProxyPostFavorHandler) SendSucceed() {
	//log.Println("SendOk")
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "",
	})
}
