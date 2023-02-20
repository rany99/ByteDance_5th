package userinfo

import (
	"ByteDance_5th/models"
	"ByteDance_5th/util/cache"
	"log"
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
	if err := models.NewUserInfoDAO().GetFriendsById(q.uid, &userList); err != nil {
		return err
	}
	log.Println(len(userList))
	for i := 0; i < len(userList); i++ {
		userList[i].IsFollow = cache.NewProxyIndexMap().GetAFollowB(q.uid, userList[i].Id)
	}
	friendList := make([]*models.Friend, len(userList))
	for i, u := range userList {
		//GetLatestMsgByUid 的 err 可能由两用户没有进行过消息通信造成，非致命，继续执行
		msg, msgType, _ := GetLatestMsgByUid(q.uid, u.Id)
		friendList[i] = &models.Friend{
			UserInfo: *u,
			Message:  msg,
			MsgType:  msgType,
		}
	}
	q.userList = friendList
	//log.Println(len(friendList))
	return nil
}

// PackData 封装数据
func (q *QueryFriendsFlow) PackData() error {
	q.FriendsResponse = &FriendsResponse{FriendList: q.userList}
	return nil
}

// GetLatestMsgByUid 通过user_id 返回最新一条聊天记录的内容以及聊天记录的类型
func GetLatestMsgByUid(fromId int64, toId int64) (string, int64, error) {
	msg, msgType, err := models.NewMessageDAO().QueryLatestMsgByUid(fromId, toId)
	if err != nil {
		return "", 0, err
	}
	return msg, msgType, nil
}
