package user_info

import (
	"ByteDance_5th/models"
	"log"
)

type Follows struct {
	UserList []*models.UserInfo `json:"user_list"`
}

type QueryFollowsFlow struct {
	uid      int64
	userList []*models.UserInfo
	*Follows
}

func QueryFollowList(uid int64) (*Follows, error) {
	return NewQueryFollowsFlow(uid).Do()
}

func NewQueryFollowsFlow(uid int64) *QueryFollowsFlow {
	return &QueryFollowsFlow{
		uid: uid,
	}
}

func (q *QueryFollowsFlow) Do() (*Follows, error) {
	if err := q.CheckJson(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.Follows, nil
}

func (q *QueryFollowsFlow) CheckJson() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

func (q *QueryFollowsFlow) GetData() error {
	var userList []*models.UserInfo
	log.Println(q.uid)
	if err := models.NewUserInfoDAO().GetFollowsById(q.uid, &userList); err != nil {
		return err
	}
	for i, _ := range userList {
		userList[i].IsFollow = true
	}
	log.Println("GetData:", len(userList))
	q.userList = userList
	log.Println("GetData:", len(q.userList))
	return nil
}

func (q *QueryFollowsFlow) PackData() error {
	log.Println("QueryFollowsFlow: PackData", len(q.userList))
	q.Follows = &Follows{UserList: q.userList}
	return nil
}
