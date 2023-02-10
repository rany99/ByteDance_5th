package video

import (
	"ByteDance_5th/models"
	"errors"
	"github.com/gin-gonic/gin"
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
	actionType string
}

func (p *ProxyPostFavorHandler) Do() {
	if err := p.ParseNum(); err != nil {
		p.SendErr(err.Error())
	}
}

// ParseNum 解析UserId
func (p *ProxyPostFavorHandler) ParseNum() error {
	RawUserid, _ := p.Get("user_id")
	userId, ok := RawUserid.(int64)
	if !ok {
		return errors.New("userid解析错误")
	}
	RawVideoId := p.Query("video_id")
	videoId, err := strconv.ParseInt(RawVideoId, 10, 64)
	if err != nil {
		return err
	}
	actionType := p.Query("action_type")
	if actionType != "1" && actionType != "2" {
		return errors.New("actionType仅限1点赞2取消")
	}
	p.videoId, p.actionType, p.userId = videoId, actionType, userId
	return nil
}

// SendErr 生成错误返回
func (p *ProxyPostFavorHandler) SendErr(msg string) {
	p.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

// SendOk 生成正确返回
func (p *ProxyPostFavorHandler) SendOk() {
	p.JSON(http.StatusOK, models.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "",
	})
}
