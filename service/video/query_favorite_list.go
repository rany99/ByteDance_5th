package video

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util/cache"
	"sync"
)

type FavoriteListResponse struct {
	Videos []*models.Video `json:"video_list"`
}

type QueryFavoriteListFlow struct {
	uidQuery  int64
	uid       int64
	videos    []*models.Video
	videoList *FavoriteListResponse
}

func QueryFavoriteList(uidQuery int64, uid int64) (*FavoriteListResponse, error) {
	return NewQueryFavoriteListFlow(uidQuery, uid).Operation()
}

func NewQueryFavoriteListFlow(uidQuery int64, uid int64) *QueryFavoriteListFlow {
	return &QueryFavoriteListFlow{
		uidQuery: uidQuery,
		uid:      uid,
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
	if err := models.NewVideoDao().QueryFavorListByUserId(q.uidQuery, &q.videos); err != nil {
		return err
	}
	p := cache.NewProxyIndexMap()
	wg := sync.WaitGroup{}
	wg.Add(len(q.videos))
	for i := range q.videos {
		go func(i int, videos *[]*models.Video, p *cache.ProxyCache) {
			var author models.UserInfo
			if err := models.NewUserInfoDAO().QueryUserInfoById((*videos)[i].UserInfoId, &author); err == nil {
				(*videos)[i].Author = author
			}
			(*videos)[i].IsFavorite = p.GetVideoFavor(q.uid, (*videos)[i].Id)
			wg.Done()
		}(i, &q.videos, p)
	}
	wg.Wait()
	return nil
}

// PackData 封装数据
func (q *QueryFavoriteListFlow) PackData() error {
	q.videoList = &FavoriteListResponse{Videos: q.videos}
	return nil
}
