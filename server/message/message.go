package message

import "ByteDance_5th/models"

type MList struct {
	Messages []*models.Message `json:"comment_list"`
}

type QueryMessageListFlow struct {
	fromId   int64
	toId     int64
	messages []*models.Message
	MList    *MList
}

func QueryMessageList(fromId int64, toId int64) (*MList, error) {
	return NewQueryMessageListFlow(fromId, toId).Operation()
}

func NewQueryMessageListFlow(fromId int64, toId int64) *QueryMessageListFlow {
	return &QueryMessageListFlow{
		fromId: fromId,
		toId:   toId,
	}
}

func (q *QueryMessageListFlow) Operation() (*MList, error) {
	if err := q.CheckJSON(); err != nil {
		return nil, err
	}
	return q.MList, nil
}

func (q *QueryMessageListFlow) GetData() error {
	if err := models.NewMessageDAO().QueryMsgListByFromIdAndToId(q.fromId, q.toId, &q.messages); err != nil {
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
