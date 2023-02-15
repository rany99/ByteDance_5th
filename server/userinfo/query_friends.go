package userinfo

import (
	"ByteDance_5th/cache"
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/common"
	"fmt"
	"log"
	"strconv"
)

type FriendsResponse struct {
	FriendList []*models.Friend `json:"user_list"`
}

type QueryFriendsFlow struct {
	uid      int64
	userList []*models.Friend
	*FriendsResponse
}

func QueryFriends(uid int64) (*FriendsResponse, error) {
	return NewQueryFriendsFlow(uid).Operation()
}

func NewQueryFriendsFlow(uid int64) *QueryFriendsFlow {
	return &QueryFriendsFlow{uid: uid}
}

func (q *QueryFriendsFlow) Operation() (*FriendsResponse, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.FriendsResponse, nil
}

// CheckJSON 校验JSON
func (q *QueryFriendsFlow) CheckJSON() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.uid); err != nil {
		return err
	}
	return nil
}

// GetData 获取朋友列表信息
func (q *QueryFriendsFlow) GetData() error {
	var userList []*models.UserInfo
	if err := models.NewUserInfoDAO().GetFansById(q.uid, &userList); err != nil {
		return err
	}
	for i := 0; i < len(userList); i++ {
		userList[i].IsFollow = cache.NewProxyIndexMap().GetAFollowB(q.uid, userList[i].Id)
	}
	friendList := make([]*models.Friend, len(userList))
	for i, u := range userList {
		friendList[i] = &models.Friend{
			UserInfo: *u,
			Avatar:   GetAvatarUrl(i),
		}
	}
	q.userList = friendList
	log.Println(len(friendList))
	return nil
}

// PackData 封装数据
func (q *QueryFriendsFlow) PackData() error {
	q.FriendsResponse = &FriendsResponse{FriendList: q.userList}
	return nil
}

// GetAvatarUrl 生成头像url
// 由于本次客户端中并没有给出相应的用于上传头像的接口，因此在public/avatar文件中预存了16张图片用作头像
func GetAvatarUrl(i int) string {
	fileName := strconv.Itoa(i%models.AvatarCnt) + ".jpg"
	var url string = fmt.Sprintf("http://%s:%d/static/avatar/%s", common.Conf.SE.IP, common.Conf.SE.Port, fileName)
	return url
}
