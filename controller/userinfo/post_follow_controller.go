package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/server/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProxyPostFollow struct {
	//uid关注/取消关注uidFollowed操作
	uid        int64
	toUserId   int64
	actionType int64
	*gin.Context
}

func PostFollowController(ctx *gin.Context) {
	NewProxyPostFollow(ctx).Operation()
}

func NewProxyPostFollow(context *gin.Context) *ProxyPostFollow {
	return &ProxyPostFollow{Context: context}
}

func (p *ProxyPostFollow) Operation() {
	if err := p.ParseJSON(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	if err := p.FollowAction(); err != nil {
		p.SendFailed(err.Error())
		return
	}
	p.SendSuccessfully()
}

func (p *ProxyPostFollow) ParseJSON() error {
	//解析uid
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New("ProxyPostFollow：uid解析错误")
	}

	//解析被关注者id
	rawToUserId := p.Query("to_user_id")
	toUserId, err := strconv.ParseInt(rawToUserId, 10, 64)
	if err != nil {
		return nil
	}

	//解析action_type
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return nil
	}
	p.uid, p.toUserId, p.actionType = uid, toUserId, actionType
	return nil
}

func (p *ProxyPostFollow) SendFailed(msg string) {
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (p *ProxyPostFollow) SendSuccessfully() {
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "",
	})
}

func (p *ProxyPostFollow) FollowAction() error {
	if stateCode := userinfo.PostFollow(p.uid, p.toUserId, p.actionType); stateCode != userinfo.NoErrOR0 {
		switch stateCode {
		case userinfo.ERRTYPE1:
			return errors.New("被关注的用户不存在")
		case userinfo.ERRTYPE2:
			return errors.New("传入ActionType只能为1或2")
		case userinfo.ERRTYPE3:
			return errors.New("不能自己关注自己")
		case userinfo.ERRTYPE4:
			return errors.New("您已关注，请勿重复关注")
		}
	}
	return nil
}
