package video

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/server/video"
	"ByteDance_5th/util"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	common.CommonResponse
	*video.FeedList
}

func FeedListController(ctx *gin.Context) {
	p := NewProxyFeedList(ctx)
	token, ok := ctx.GetQuery("token")
	//log.Println(token)
	if !ok || token == "" {
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

// ProxyFeedList 代理层
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
	timeMs, err := strconv.ParseInt(timeStamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, timeMs*1e6)
	}
	//log.Println("无token时间戳======>>>>>", latestTime)
	list, err := video.QueryFeedList(0, latestTime)
	if err != nil {
		return err
	}
	p.GetFeedListSuccessfully(list)
	return nil
}

// DoWithToken 登陆状态下的视频推送
func (p *ProxyFeedList) DoWithToken(token string) error {
	if claim, ok := util.DecodeToken(token); ok {
		if time.Now().Unix() > claim.ExpiresAt {
			return errors.New(errortype.TokenOutDateErr)
		}
		timeStamp := p.Query("latest_time")
		var latestTime time.Time
		iniTime, err := strconv.ParseInt(timeStamp, 10, 24)
		if err == nil {
			latestTime = time.Unix(0, iniTime*1e6)
		}
		list, err := video.QueryFeedList(claim.UserId, latestTime)
		if err != nil {
			return err
		}
		p.GetFeedListSuccessfully(list)
		return nil
	}
	return errors.New(errortype.ParseTokenErr)
}

// GetFeedListSuccessfully 获取视频流成功
func (p *ProxyFeedList) GetFeedListSuccessfully(feedList *video.FeedList) {
	//if feedList == nil {
	//	log.Println("GetFeedListSuccessfully")
	//	return
	//}
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		FeedList: feedList,
	})

}

// GetFeedListFailed 获取视频流失败
func (p *ProxyFeedList) GetFeedListFailed(msg string) {
	p.JSON(http.StatusOK, FeedResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}
