package video

import (
	"ByteDance_5th/models"
)

type FavoriteListResponse struct {
	Videos []*models.Video `json:"video_list"`
}

type QueryFavoriteListFlow struct {
	uid       int64
	videos    []*models.Video
	videoList *FavoriteListResponse
}

func QueryFavoriteList(uid int64) (*FavoriteListResponse, error) {
	return NewQueryFavoriteListFlow(uid).Operation()
}

func NewQueryFavoriteListFlow(uid int64) *QueryFavoriteListFlow {
	return &QueryFavoriteListFlow{
		uid: uid,
	}
}

func (q *QueryFavoriteListFlow) Operation() (*FavoriteListResponse, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.videoList, nil
}

// CheckJSON 校验数据合法性
func (q *QueryFavoriteListFlow) CheckJSON() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

// GetData 调用DAO层
func (q *QueryFavoriteListFlow) GetData() error {
	if err := models.NewVideoDao().QueryFavorListByUserId(q.uid, &q.videos); err != nil {
		return err
	}
	for i := 0; i < len(q.videos); i++ {
		var author models.UserInfo
		if err := models.NewUserInfoDAO().QueryUserInfoById(q.videos[i].UserInfoId, &author); err != nil {
			q.videos[i].Author = author
		}
		q.videos[i].IsFavorite = true
	}
	return nil
}

// PackData 封装数据
func (q *QueryFavoriteListFlow) PackData() error {
	q.videoList = &FavoriteListResponse{Videos: q.videos}
	return nil
}