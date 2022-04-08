package model

import (
	"encoding/json"
	"fmt"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"time"
)

const (
	MaintainMsgType = iota
	CommonMsgType
	HeartbeatMsgType
	AckMsgType
)

type Message struct {
	Id          int64  `json:"id" gorm:"column:id"`
	Time        int64  `json:"time" gorm:"column:time"`
	UserId      int    `json:"user_id" gorm:"column:user_id"`
	RoomId      string `json:"room_id" gorm:"column:room_id"` // 64B
	MessageType uint8  `json:"message_type" gorm:"message_type"`
	HasRead     bool   `json:"had_read" gorm:"column:has_read"`
	HasSend     bool   `json:"has_send" gorm:"column:has_send"`
	Text        string `json:"text" gorm:"column:text"`
}

func (_ Message) TableName() string {
	return "message"
}

func NewMessage(userId int, roomId string, msgType uint8, text string) Message {
	return Message{
		Id:          common.GenerateID(),
		Time:        time.Now().UnixNano(),
		UserId:      userId,
		RoomId:      roomId,
		MessageType: msgType,
		Text:        text,
	}
}

func (m Message) String() string {
	//TODO: better to_string
	return fmt.Sprintf("%#v", m)
}

func (m Message) MarshalBinary() ([]byte, error) {
	//TODO: better marshal binary
	marshal, err := json.Marshal(m)
	if err != nil {
		log.LOG.Panic(err)
		return nil, err
	}
	return marshal, nil
}

func ParseMessage(raw []byte) (Message, error) {
	var ret Message
	err := json.Unmarshal(raw, &ret)
	if err != nil {
		log.LOG.Errorf("fail to parse message, Err=%+v", err)
	}
	return ret, err
}

func NewNotifyMessage(text string, targetRoom string) *Message {
	msg := NewMessage(0, targetRoom, MaintainMsgType, text)
	return &msg
}
