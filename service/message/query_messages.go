package message

import (
	"ByteDance_5th/models"
)

type MessagesResponse struct {
	Messages []*models.Message `json:"message_list"`
}

type QueryMessageListFlow struct {
	fromId           int64
	toId             int64
	preMsgTime       int64
	messages         []*models.Message
	messagesResponse *MessagesResponse
}

func QueryMessageList(fromId int64, toId int64, preMsgTime int64) (*MessagesResponse, error) {
	return NewQueryMessageListFlow(fromId, toId, preMsgTime).Operation()
}

func NewQueryMessageListFlow(fromId int64, toId int64, preMsgTime int64) *QueryMessageListFlow {
	return &QueryMessageListFlow{
		fromId:     fromId,
		toId:       toId,
		preMsgTime: preMsgTime,
	}
}

func (q *QueryMessageListFlow) Operation() (*MessagesResponse, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	if err := q.GetData(); err != nil {
		return nil, err
	}
	if err := q.PackData(); err != nil {
		return nil, err
	}
	return q.messagesResponse, nil
}

func (q *QueryMessageListFlow) GetData() error {
	if err := models.NewMessageDAO().QueryMsgListByFromIdAndToId(q.fromId, q.toId, &q.messages, q.preMsgTime); err != nil {
		return err
	}
	return nil
}

func (q *QueryMessageListFlow) CheckJSON() error {
	//判断发出者是否存在
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.fromId); err != nil {
		return err
	}
	//判断接收者是否存在
	if err := models.NewUserInfoDAO().IsUserInfoExist(q.toId); err != nil {
		return err
	}
	return nil
}

func (q *QueryMessageListFlow) PackData() error {
	q.messagesResponse = &MessagesResponse{Messages: q.messages}
	return nil
}
