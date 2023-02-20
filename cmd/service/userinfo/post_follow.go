package userinfo

import (
	"ByteDance_5th/cmd/models"
	"ByteDance_5th/util/cache"
)

const (
	NoErrOR0 = 0 //没有错误
	ERRTYPE1 = 1 //被关注的用户不存在
	ERRTYPE2 = 2 //ActionType数据为1关注2取消以外的其他选项
	ERRTYPE3 = 3 //不能自己关注自己
	ERRTYPE4 = 4 //数据库操作失败
)

type PostFollowFlow struct {
	uid        int64
	toUserId   int64
	actionType int64
}

func PostFollow(uid, toUserId, actionType int64) int {
	return NewPostFollowFlow(uid, toUserId, actionType).Do()
}

func NewPostFollowFlow(uid int64, toUserId int64, actionType int64) *PostFollowFlow {
	return &PostFollowFlow{
		uid:        uid,
		toUserId:   toUserId,
		actionType: actionType,
	}
}

func (p *PostFollowFlow) Do() int {
	stateCode := 0
	if stateCode = p.CheckJson(); stateCode != NoErrOR0 {
		return stateCode
	}
	if stateCode = p.DoFollowAction(); stateCode != NoErrOR0 {
		return stateCode
	}
	return stateCode
}

func (p *PostFollowFlow) CheckJson() int {
	if err := models.NewUserInfoDAO().IsUserInfoExist(p.uid); err != nil {
		return ERRTYPE1
	}
	//log.Println(p.actionType)
	if p.actionType != 1 && p.actionType != 2 {
		return ERRTYPE2
	}
	if p.uid == p.toUserId {
		return ERRTYPE3
	}
	return NoErrOR0
}

func (p *PostFollowFlow) DoFollowAction() int {
	userDAO := models.NewUserInfoDAO()
	stateCode := NoErrOR0
	switch p.actionType {
	case 1:
		if err := userDAO.AFollowB(p.uid, p.toUserId); err != nil {
			stateCode = ERRTYPE4
		}
		cache.NewProxyIndexMap().SetAFollowB(p.uid, p.toUserId, true)
	case 2:
		if err := userDAO.NoAFollowB(p.uid, p.toUserId); err != nil {
			stateCode = ERRTYPE4
		}
		cache.NewProxyIndexMap().SetAFollowB(p.uid, p.toUserId, false)
	default:
		stateCode = ERRTYPE2
	}
	return stateCode
}
