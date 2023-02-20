package message

import (
	models2 "ByteDance_5th/cmd/models"
)

type MessagesResponse struct {
	Messages []*models2.Message `json:"message_list"`
}

type QueryMessageListFlow struct {
	fromId   int64
	toId     int64
	messages []*models2.Message
	MList    *MessagesResponse
}

func QueryMessageList(fromId int64, toId int64) (*MessagesResponse, error) {
	return NewQueryMessageListFlow(fromId, toId).Operation()
}

func NewQueryMessageListFlow(fromId int64, toId int64) *QueryMessageListFlow {
	return &QueryMessageListFlow{
		fromId: fromId,
		toId:   toId,
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
	return q.MList, nil
}

func (q *QueryMessageListFlow) GetData() error {
	if err := models2.NewMessageDAO().QueryMsgListByFromIdAndToId(q.fromId, q.toId, &q.messages); err != nil {
		return err
	}
	return nil
}

func (q *QueryMessageListFlow) CheckJSON() error {
	//判断发出者是否存在
	if err := models2.NewUserInfoDAO().IsUserInfoExist(q.fromId); err != nil {
		return err
	}
	//判断接收者是否存在
	if err := models2.NewUserInfoDAO().IsUserInfoExist(q.toId); err != nil {
		return err
	}
	return nil
}

func (q *QueryMessageListFlow) PackData() error {
	q.MList = &MessagesResponse{Messages: q.messages}
	return nil
}
