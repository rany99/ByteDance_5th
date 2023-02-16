package models

import (
	"log"
	"testing"
	"time"
)

func TestMessage(t *testing.T) {
	InitDB()
	m := NewMessageDAO()
	for i := 0; i < 10; i++ {
		msg := Message{
			FromUserId: 1,
			ToUserId:   2,
			Content:    "我叫叶茂",
			CreateTime: time.Now().Unix(),
		}
		if err := m.CreateMessage(&msg); err != nil {
			log.Println(err.Error())
		}
	}
	//
	//list := []*Message{}
	//m.QueryMsgListByFromIdAndToId(1, 2, &list)
	//fmt.Println(len(list))
	//for _, msg := range list {
	//	fmt.Println(msg.Content, " ", msg.Optional)
	//}
}
