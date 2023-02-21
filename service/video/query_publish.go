package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util/cache"
	"errors"
	"sync"
)

type PublishListResponse struct {
	Videos []*models.Video `json:"video_list,omitempty"`
}

type QueryPublishListByUidFlow struct {
	uidQuery  int64
	userId    int64 //当前用户id
	videos    []*models.Video
	videoList *PublishListResponse
}

// QueryPublishListByUid 通过UID返回
func QueryPublishListByUid(uidQuery int64, userid int64) (*PublishListResponse, error) {
	return NewQueryUserVideoListByUid(uidQuery, userid).Operation()
}

func (q *QueryPublishListByUidFlow) Operation() (*PublishListResponse, error) {
	if err := q.IsUidExist(); err != nil {
		return nil, errors.New(errortype.UserNoExistErr)
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

func (q *QueryPublishListByUidFlow) IsUidExist() error {
	return models.NewUserInfoDAO().IsUserInfoExist(q.userId)
}

func NewQueryUserVideoListByUid(uidQuery int64, userid int64) *QueryPublishListByUidFlow {
	return &QueryPublishListByUidFlow{
		uidQuery: uidQuery,
		userId:   userid}
}

// PackData 封装数据
func (q *QueryPublishListByUidFlow) PackData() error {
	var userInfo models.UserInfo
	if err := models.NewUserInfoDAO().QueryUserInfoById(q.uidQuery, &userInfo); err != nil {
		return err
	}
	if err := models.NewVideoDao().QueryVideoListByUserId(q.userId, &q.videos); err != nil {
		return err
	}

	p := cache.NewProxyIndexMap()

	wg := sync.WaitGroup{}
	wg.Add(len(q.videos))

	for i := range q.videos {
		go func(i int, userInfo *models.UserInfo) {
			q.videos[i].Author = *userInfo
			q.videos[i].IsFavorite = p.GetVideoFavor(q.userId, q.videos[i].Id)
			wg.Done()
		}(i, &userInfo)
	}
	wg.Wait()
	//log.Println("PackData:", len(q.videos))
	q.videoList = &PublishListResponse{Videos: q.videos}
	return nil
}
