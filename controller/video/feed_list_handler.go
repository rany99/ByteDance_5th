package video

import (
	"ByteDance_5th/middle"
	"ByteDance_5th/models"
	"ByteDance_5th/server/video"
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	models.CommonResponse
	*video.FeedList
}

func FeedListHandler(ctx *gin.Context) {
	p := NewProxyFeedList(ctx)
	token, ok := ctx.GetQuery("token")
	if !ok || token == "" {
		log.Println("无token")
		err := p.DoWithoutToken()
		if err != nil {
			p.GetFeedListFailed(err.Error())
		}
		return
	}

	err := p.DoWithToken(token)
	if err != nil {
		p.GetFeedListFailed(err.Error())
	}
}

// ProxyFeedList f防火层
type ProxyFeedList struct {
	*gin.Context
}

func NewProxyFeedList(ctx *gin.Context) *ProxyFeedList {
	return &ProxyFeedList{ctx}
}

// DoWithoutToken 未登陆状态下的视频推送
func (p *ProxyFeedList) DoWithoutToken() error {
	var latestTime time.Time
	timeStamp := p.Query("latest_time")
	log.Println("timeStamp:", timeStamp)
	timeMs, err := strconv.ParseInt(timeStamp, 10, 64)
	log.Println("timeMs:", timeMs)
	if err != nil {
		latestTime = time.Unix(0, timeMs*1e6)
	}
	list, err := video.QueryFeedList(0, latestTime)
	if err != nil {
		return err
	}
	p.GetFeedListSuccessfully(list)
	return nil
}

// DoWithToken 登陆状态下的视频推送
func (p *ProxyFeedList) DoWithToken(token string) error {
	if claim, ok := middle.DecodeToken(token); ok {
		if time.Now().Unix() > claim.ExpiresAt {
			return errors.New("token已经过期")
		}
		timeStamp := p.Query("latest_time")
		var latestTime time.Time
		iniTime, err := strconv.ParseInt(timeStamp, 10, 24)
		if err != nil {
			latestTime = time.Unix(0, iniTime*1e6)
		}
		list, err := video.QueryFeedList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		p.GetFeedListSuccessfully(list)
		return nil
	}
	return errors.New("token解析错误")
}

// GetFeedListSuccessfully 获取视频流成功
func (p *ProxyFeedList) GetFeedListSuccessfully(feedList *video.FeedList) {
	if feedList == nil {
		log.Println("GetFeedListSuccessfully：视频流指针为空")
	}
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		FeedList: feedList,
	})

}

// GetFeedListFailed 获取视频流失败
func (p *ProxyFeedList) GetFeedListFailed(msg string) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: models.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}
