package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util/cache"
	"log"
	"sync"
)

type FollowsResponse struct {
	UserList []*models.UserInfo `json:"user_list"`
}

type QueryFollowsFlow struct {
	uidQuery int64
	uid      int64
	userList []*models.UserInfo
	*FollowsResponse
}

func QueryFollowList(uidQuery int64, uid int64) (*FollowsResponse, error) {
	return NewQueryFollowsFlow(uidQuery, uid).Operation()
}

func NewQueryFollowsFlow(uidQuery int64, uid int64) *QueryFollowsFlow {
	return &QueryFollowsFlow{
		uidQuery: uidQuery,
		uid:      uid,
	}
}

func (q *QueryFollowsFlow) Operation() (*FollowsResponse, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.FollowsResponse, nil
}

func (q *QueryFollowsFlow) CheckJSON() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uidQuery); err != nil {
		return err
	}
	return nil
}

func (q *QueryFollowsFlow) GetData() error {
	log.Println(q.uid)
	log.Println(q.uidQuery)
	var userList []*models.UserInfo
	if err := models.NewUserInfoDAO().GetFollowsById(q.uidQuery, &userList); err != nil {
		return err
	}
	p := cache.NewProxyIndexMap()
	wg := sync.WaitGroup{}
	wg.Add(len(userList))
	for i := range userList {
		go func(i int, p *cache.ProxyCache) {
			userList[i].IsFollow = p.GetAFollowB(q.uid, userList[i].Id)
			wg.Done()
		}(i, p)
	}
	wg.Wait()
	q.userList = userList
	return nil
}

func (q *QueryFollowsFlow) PackData() error {
	q.FollowsResponse = &FollowsResponse{UserList: q.userList}
	return nil
}
