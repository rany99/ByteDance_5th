package message

import (
	models2 "ByteDance_5th/cmd/models"
	"ByteDance_5th/pkg/errortype"
	"errors"
	"time"
)

// PostMessageFlow 网络层模型
type PostMessageFlow struct {
	fromId     int64
	toId       int64
	actionType int64
	content    string
}

// PostMessage 发送信息
func PostMessage(fromId int64, toId int64, actionType int64, content string) error {
	return NewPostMessageFlow(fromId, toId, actionType, content).Operation()
}

// NewPostMessageFlow 新的网络层模型
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

// CheckJSON 校验JSON解析得到的参数信息
func (p *PostMessageFlow) CheckJSON() error {
	// 发出者是否存在
	if err := models2.NewUserInfoDAO().IsUserInfoExist(p.fromId); err != nil {
		return errors.New(errortype.FromUserNoExistErr)
	}
	// 接收者是否存在
	if err := models2.NewUserInfoDAO().IsUserInfoExist(p.toId); err != nil {
		return errors.New(errortype.ToUserNoExistErr)
	}
	// 消息是否为空
	if p.content == "" {
		return errors.New(errortype.EmptyMsgErr)
	}
	return nil
}

// GetData 调用DAL操作
func (p *PostMessageFlow) GetData() error {
	message := models2.Message{
		UserInfoId: p.fromId,
		ToUserId:   p.toId,
		Content:    p.content,
		CreateTime: time.Now().Unix(),
	}
	return models2.NewMessageDAO().CreateMessage(&message)
}
