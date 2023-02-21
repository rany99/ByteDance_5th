package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util/cache"
	"sync"
)

type FansResponse struct {
	UserList []*models.UserInfo `json:"user_list"`
}

type QueryFansFlow struct {
	uidQuery int64
	uid      int64
	userList []*models.UserInfo
	*FansResponse
}

func QueryFans(uidQuery int64, uid int64) (*FansResponse, error) {
	return NewQueryFansFlow(uidQuery, uid).Operation()
}

func NewQueryFansFlow(uidQuery int64, uid int64) *QueryFansFlow {
	return &QueryFansFlow{
		uidQuery: uidQuery,
		uid:      uid,
	}
}

func (q *QueryFansFlow) Operation() (*FansResponse, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.FansResponse, nil
}

func (q *QueryFansFlow) CheckJSON() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

func (q *QueryFansFlow) GetData() error {
	if err := models.NewUserInfoDAO().GetFansById(q.uidQuery, &q.userList); err != nil {
		return err
	}
	p := cache.NewProxyIndexMap()
	wg := sync.WaitGroup{}
	wg.Add(len(q.userList))
	for i := range q.userList {
		go func(i int, p *cache.ProxyCache) {
			q.userList[i].IsFollow = cache.NewProxyIndexMap().GetAFollowB(q.uid, q.userList[i].Id)
			wg.Done()
		}(i, p)
	}
	wg.Wait()
	return nil
}

func (q *QueryFansFlow) PackData() error {
	q.FansResponse = &FansResponse{UserList: q.userList}
	return nil
}
