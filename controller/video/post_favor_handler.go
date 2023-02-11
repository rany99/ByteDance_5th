package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/server/video"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func PostFavorHandler(ctx *gin.Context) {
	NewPostFavorHandler(ctx).Do()
}

func NewPostFavorHandler(ctx *gin.Context) *ProxyPostFavorHandler {
	return &ProxyPostFavorHandler{Context: ctx}
}

type ProxyPostFavorHandler struct {
	*gin.Context
	userId     int64
	videoId    int64
	actionType int64
}

func (p *ProxyPostFavorHandler) Do() {
	//解码
	if err := p.ParseNum(); err != nil {
		p.SendErr(err.Error())
		return
	}
	//sever层执行点赞操作
	if err := video.PostFavor(p.userId, p.videoId, p.actionType); err != nil {
		p.SendErr(err.Error())
		return
	}
	//点赞成功
	p.SendOk()
}

// ParseNum 解析UserId
func (p *ProxyPostFavorHandler) ParseNum() error {
	RawUserid, _ := p.Get("user_id")
	userId, ok := RawUserid.(int64)
	if !ok {
		return errors.New("userid解析错误")
	}
	log.Println("userId解析成功：", userId)
	RawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(RawVideoId, 10, 64)
	log.Println("videoId解析成功：", videoId)
	if err != nil {
		return err
	}
	RawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(RawActionType, 10, 64)
	if err != nil {
		return err
	}
	//log.Println("actionType解析成功：", actionType)
	if actionType != 1 && actionType != 2 {
		return errors.New("actionType仅限1点赞2取消")
	}
	p.videoId, p.actionType, p.userId = videoId, actionType, userId
	return nil
}

// SendErr 生成错误返回
func (p *ProxyPostFavorHandler) SendErr(msg string) {
	log.Println("SendErr")
	p.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// SendOk 生成正确返回
func (p *ProxyPostFavorHandler) SendOk() {
	log.Println("SendOk")
	p.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "",
	})
}
