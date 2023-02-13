package userinfo

import (
	"ByteDance_5th/cache"
	"ByteDance_5th/models"
)

type FansResponse struct {
	UserList []*models.UserInfo `json:"user_list"`
}

type QueryFansFlow struct {
	uid      int64
	userList []*models.UserInfo
	*FansResponse
}

func QueryFans(uid int64) (*FansResponse, error) {
	return NewQueryFansFlow(uid).Do()
}

func NewQueryFansFlow(uid int64) *QueryFansFlow {
	return &QueryFansFlow{uid: uid}
}

func (q *QueryFansFlow) Do() (*FansResponse, error) {
	if err := q.CheckJson(); err != nil {
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

func (q *QueryFansFlow) CheckJson() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

func (q *QueryFansFlow) GetData() error {
	var userList []*models.UserInfo
	if err := models.NewUserInfoDAO().GetFansById(q.uid, &userList); err != nil {
		return err
	}
	//log.Println("server层：列表长度", len(userList))
	for i := 0; i < len(userList); i++ {
		userList[i].IsFollow = cache.NewProxyIndexMap().GetAFollowB(q.uid, userList[i].Id)
	}
	q.userList = userList
	return nil
}

func (q *QueryFansFlow) PackData() error {
	q.FansResponse = &FansResponse{UserList: q.userList}
	return nil
}
