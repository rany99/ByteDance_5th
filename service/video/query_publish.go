package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/util/cache"
	"errors"
)

type PublishListResponse struct {
	Videos []*models.Video `json:"video_list,omitempty"`
}

type QueryPublishListByUidFlow struct {
	userId    int64 //当前用户id
	videos    []*models.Video
	videoList *PublishListResponse
}

// QueryPublishListByUid 通过UID返回
func QueryPublishListByUid(userid int64) (*PublishListResponse, error) {
	return NewQueryUserVideoListByUid(userid).Operation()
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

func NewQueryUserVideoListByUid(userid int64) *QueryPublishListByUidFlow {
	return &QueryPublishListByUidFlow{userId: userid}
}

// PackData 封装数据
func (q *QueryPublishListByUidFlow) PackData() error {
	if err := models.NewVideoDao().QueryVideoListByUserId(q.userId, &q.videos); err != nil {
		return err
	}
	var userInfo models.UserInfo
	if err := models.NewUserInfoDAO().QueryUserInfoById(q.userId, &userInfo); err != nil {
		return err
	}
	p := cache.NewProxyIndexMap()

	for i := range q.videos {
		q.videos[i].Author = userInfo
		q.videos[i].IsFavorite = p.GetVideoFavor(q.userId, q.videos[i].Id)
	}
	//log.Println("PackData:", len(q.videos))
	q.videoList = &PublishListResponse{Videos: q.videos}
	return nil
}
