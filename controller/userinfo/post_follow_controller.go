package userinfo

import (
	"ByteDance_5th/pkg/common"
	"ByteDance_5th/pkg/errortype"
	"ByteDance_5th/service/userinfo"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ProxyPostFollow struct {
	//uid关注/取消关注toUserId操作
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
	p.SendSucceed()
}

// ParseJSON 校验传入的JSON信息
func (p *ProxyPostFollow) ParseJSON() error {
	//解析uid
	rawUid, _ := p.Get("user_id")
	uid, ok := rawUid.(int64)
	if !ok {
		return errors.New(errortype.ParseUserIdErr)
	}

	//解析被关注者id
	rawToUserId := p.Query("to_user_id")
	toUserId, err := strconv.ParseInt(rawToUserId, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseToUserIdErr)
	}

	//解析action_type
	rawActionType := p.Query("action_type")
	actionType, err := strconv.ParseInt(rawActionType, 10, 64)
	if err != nil {
		return errors.New(errortype.ParseActionTypeErr)
	}

	//填入代理层
	p.uid, p.toUserId, p.actionType = uid, toUserId, actionType
	return nil
}

func (p *ProxyPostFollow) SendFailed(msg string) {
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 1,
		StatusMsg:  msg,
	})
}

func (p *ProxyPostFollow) SendSucceed() {
	p.JSON(http.StatusOK, common.CommonResponse{
		StatusCode: 0,
		StatusMsg:  "关注成功",
	})
}

func (p *ProxyPostFollow) FollowAction() error {
	if stateCode := userinfo.PostFollow(p.uid, p.toUserId, p.actionType); stateCode != userinfo.NoErrOR0 {
		switch stateCode {
		case userinfo.ERRTYPE1:
			return errors.New(errortype.FollowUserNoExistErr)
		case userinfo.ERRTYPE2:
			return errors.New(errortype.PostFollowActionTypeErr)
		case userinfo.ERRTYPE3:
			return errors.New(errortype.CantFollowSelfErr)
		case userinfo.ERRTYPE4:
			return errors.New(errortype.FollowAgainErr)
		}
	}
	return nil
}
