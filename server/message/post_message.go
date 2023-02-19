package message

import (
	"ByteDance_5th/models"
	"ByteDance_5th/pkg/constantval"
	"ByteDance_5th/pkg/errortype"
	"errors"
)

type PostMessageFlow struct {
	fromId     int64
	toId       int64
	actionType int64
	content    string
}

func PostMessage(fromId int64, toId int64, actionType int64, content string) error {
	return NewPostMessageFlow(fromId, toId, actionType, content).Operation()
}

func NewPostMessageFlow(fromId, toId, actionType int64, content string) *PostMessageFlow {
	return &PostMessageFlow{
		fromId:     fromId,
		toId:       toId,
		actionType: actionType,
		content:    content,
	}
}

func (p *PostMessageFlow) Operation() error {
	if err := p.CheckJSON(); err != nil {
		return err
	}
	if err := p.GetData(); err != nil {
		return err
	}
	return nil
}

func (p *PostMessageFlow) CheckJSON() error {
	if err := models.NewUserInfoDAO().IsUserInfoExist(p.fromId); err != nil {
		return errors.New(errortype.FromUserNoExistErr)
	}
	if err := models.NewUserInfoDAO().IsUserInfoExist(p.toId); err != nil {
		return errors.New(errortype.ToUserNoExistErr)
	}
	if p.actionType != constantval.SendMsgActionType {
		return errors.New(errortype.PostMsgActionTypeErr)
	}
	return nil
}

func (p *PostMessageFlow) GetData() error {
	message := models.Message{
		UserInfoId: p.fromId,
		ToUserId:   p.toId,
		Content:    p.content,
	}
	return models.NewMessageDAO().CreateMessage(&message)
}
