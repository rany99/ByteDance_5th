package video

import (
	"ByteDance_5th/models"
)

type FavoriteList struct {
	Videos []*models.Video `json:"video_list"`
}

type QueryFavoriteListFlow struct {
	uid       int64
	videos    []*models.Video
	videoList *FavoriteList
}

func QueryFavoriteList(uid int64) (*FavoriteList, error) {
	return NewQueryFavoriteListFlow(uid).Operation()
}

func NewQueryFavoriteListFlow(uid int64) *QueryFavoriteListFlow {
	return &QueryFavoriteListFlow{
		uid: uid,
	}
}

func (q *QueryFavoriteListFlow) Operation() (*FavoriteList, error) {
	if err := q.CheckJson(); err != nil {
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

func (q *QueryFavoriteListFlow) CheckJson() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

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

func (q *QueryFavoriteListFlow) PackData() error {
	//log.Println("QueryFavoriteListFlow->PackData:", len(q.videos))
	//log.Println(q.videos[0].PlayUrl)
	q.videoList = &FavoriteList{Videos: q.videos}
	return nil
}
