package video

import (
	"ByteDance_5th/cmd/service/video"
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type QueryFeedResponse struct {
	common.CommonResponse
	*video.FeedListResponse
}

func QueryFeedListController(ctx *gin.Context) {

	// 解析token token非必须
	token, ok := ctx.GetQuery("token")

	// 无token状态下的视频推送
	if !ok || token == "" {
		if err := DoWithoutToken(ctx); err != nil {
			QueryFeedListFailed(ctx, err.Error())
			return
		}
		// 无状态下推送成功，返回
		return
	}

	// 有token下的视频推送
	if err := DoWithToken(ctx, token); err != nil {
		QueryFeedListFailed(ctx, err.Error())
		return
	}
}

// DoWithoutToken 未登陆状态下的视频推送
func DoWithoutToken(ctx *gin.Context) error {
	// 解析时间戳
	var latestTime time.Time
	timeStamp := ctx.Query("latest_time")
	timeMs, err := strconv.ParseInt(timeStamp, 10, 64)
	if err == nil {
		latestTime = time.Unix(0, timeMs*1e6)
	}
	//log.Println("无token下的时间戳", latestTime)

	// 调用service层
	feedListResponse, err := video.QueryFeedList(0, latestTime)
	if err != nil {
		return err
	}

	// 获取视频流成功
	QueryFeedListSucceed(ctx, feedListResponse)
	return nil
}

// DoWithToken 登陆状态下的视频推送
func DoWithToken(ctx *gin.Context, token string) error {
	if claim, ok := util.DecodeToken(token); ok {

		// 判断token是否过期
		if time.Now().Unix() > claim.ExpiresAt {
			return errors.New(errortype.TokenOutDateErr)
		}

		// 解析时间戳
		timeStamp := ctx.Query("latest_time")
		var latestTime time.Time
		iniTime, err := strconv.ParseInt(timeStamp, 10, 24)
		if err == nil {
			latestTime = time.Unix(0, iniTime*1e6)
		}

		// 调用service层
		feedListResponse, err := video.QueryFeedList(claim.UserId, latestTime)
		if err != nil {
			return err
		}

		// 获取视频流推送成功
		QueryFeedListSucceed(ctx, feedListResponse)
		return nil
	}

	// token 解析错误
	return errors.New(errortype.ParseTokenErr)
}

// QueryFeedListSucceed 获取视频流成功
func QueryFeedListSucceed(ctx *gin.Context, feedListResponse *video.FeedListResponse) {
	ctx.JSON(http.StatusOK, QueryFeedResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 0,
			StatusMsg:  "",
		},
		FeedListResponse: feedListResponse,
	})

}

// QueryFeedListFailed 获取视频流失败
func QueryFeedListFailed(ctx *gin.Context, msg string) {
	ctx.JSON(http.StatusOK, QueryFeedResponse{
		CommonResponse: common.CommonResponse{
			StatusCode: 1,
			StatusMsg:  msg,
		},
	})
}
