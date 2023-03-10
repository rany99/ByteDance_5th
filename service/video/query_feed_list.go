package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util/cache"
	"errors"
	"sync"
	"time"
)

const Limit int = 30

type FeedListResponse struct {
	List     []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

type QueryFeedListFlow struct {
	userid     int64           //用户ID
	latestTime time.Time       //申请时间
	videos     []*models.Video //返回视频列表，长度最长为Limit = 30
	nextTime   int64           //下一次申请时间
	feedVideo  *FeedListResponse
}

func QueryFeedList(userid int64, latestTime time.Time) (*FeedListResponse, error) {
	return NewQueryFeedListFlow(userid, latestTime).Do()
}

func NewQueryFeedListFlow(userId int64, latestTime time.Time) *QueryFeedListFlow {
	return &QueryFeedListFlow{
		userid:     userId,
		latestTime: latestTime,
	}
}

func (q *QueryFeedListFlow) Do() (*FeedListResponse, error) {

	q.IsAlreadyLogin()

	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.feedVideo, nil
}

// IsAlreadyLogin 无论是否登陆都进行视频推送
func (q *QueryFeedListFlow) IsAlreadyLogin() {
	//userid大于零表示已经登陆
	if q.userid > 0 {
		//q.videos = SelectVideosByUserId(q.userid)
	}
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
	}
}

// GetData 获取数据
func (q *QueryFeedListFlow) GetData() error {
	if err := models.NewVideoDao().QueryVideoListByLastTimeAndLimit(q.latestTime, Limit, &q.videos); err != nil {
		return err
	}
	latestTime, _ := FillVideos(q.userid, &q.videos)

	if latestTime != nil {
		q.nextTime = (*latestTime).UnixNano() / 1e6
		return nil
	}

	q.nextTime = time.Now().Unix() / 1e6
	return nil
}

// PackData 封装数据
func (q *QueryFeedListFlow) PackData() error {
	q.feedVideo = &FeedListResponse{
		List:     q.videos,
		NextTime: q.nextTime,
	}
	return nil
}

// FillVideos 更新视频作者信息
func FillVideos(userid int64, videos *[]*models.Video) (*time.Time, error) {
	videosLen := len(*videos)
	if videos == nil || videosLen == 0 {
		return nil, errors.New("FillVideos" + errortype.PointerIsNilErr)
	}

	//dao := models.NewUserInfoDAO()
	p := cache.NewProxyIndexMap()
	latestTime := (*videos)[videosLen-1].CreatedAt

	wg := sync.WaitGroup{}
	wg.Add(videosLen)

	// 依据列表长度开启若干进程填入关注与喜欢等信息
	for i := 0; i < videosLen; i++ {
		go func(i int, p *cache.ProxyCache, videos *[]*models.Video) {
			var author models.UserInfo
			if err := models.NewUserInfoDAO().QueryUserInfoById((*videos)[i].UserInfoId, &author); err != nil {
				return
			}
			author.IsFollow = p.GetAFollowB(userid, author.Id)
			(*videos)[i].Author = author
			if userid > 0 {
				(*videos)[i].IsFavorite = p.GetVideoFavor(userid, (*videos)[i].Id)
			}
			wg.Done()
		}(i, p, videos)
	}
	wg.Wait()

	return &latestTime, nil
}
