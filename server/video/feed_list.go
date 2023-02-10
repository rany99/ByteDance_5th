package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util"
	"log"
	"time"
)

const Limit int = 30

type FeedList struct {
	List     []*models.Video `json:"video_list,omitempty"`
	NextTime int64           `json:"next_time,omitempty"`
}

type QueryFeedListFlow struct {
	userid     int64           //用户ID
	latestTime time.Time       //申请时间
	videos     []*models.Video //返回视频列表，长度最长为Limit = 30
	nextTime   int64           //下一次申请时间
	feedVideo  *FeedList
}

func QueryFeedList(userid int64, latestTime time.Time) (*FeedList, error) {
	return NewQueryFeedListFlow(userid, latestTime).Do()
}

func NewQueryFeedListFlow(userId int64, latestTime time.Time) *QueryFeedListFlow {
	return &QueryFeedListFlow{
		userid:     userId,
		latestTime: latestTime,
	}
}

func (q *QueryFeedListFlow) Do() (list *FeedList, err error) {
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
		//
	}
	log.Println(time.Now())
	q.latestTime = time.Now()
	if q.latestTime.IsZero() {
		q.latestTime = time.Now()
		log.Println("latestTime->time.Now", q.latestTime)
	}
}

// GetData 获取数据
func (q *QueryFeedListFlow) GetData() error {
	if err := models.NewVideoDao().QueryVideoListByLastTimeAndLimit(q.latestTime, Limit, &q.videos); err != nil {
		return err
	}
	latestTime, _ := util.FillVideos(q.userid, &q.videos)

	if latestTime != nil {
		q.nextTime = (*latestTime).UnixNano() / 1e6
		return nil
	}

	q.nextTime = time.Now().Unix() / 1e6
	return nil
}

// PackData 封装数据
func (q *QueryFeedListFlow) PackData() error {
	log.Println("listLen", len(q.videos))
	q.feedVideo = &FeedList{
		List:     q.videos,
		NextTime: q.nextTime,
	}
	return nil
}
