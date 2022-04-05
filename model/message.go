package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const (
	MENTAIN_MSG_TYPE = iota
	COMMON_MSG_TYPE
	HEARTBEAT_MSG_TYPE
)

type Message struct {
	Time        int64  `json:"time" gorm:"column:time"`
	UserId      int64  `json:"user_id" gorm:"column:user_id"`
	RoomId      string `json:"room_id"` // 64B
	MessageType uint8  `json:"message_type"`
	Text        string `json:"text"`
}

func (_ Message) TableName() string {
	return "message"
}

func (m Message) String() string {
	//TODO: better to_string
	return fmt.Sprintf("%d,%d,%s,%d,%s", m.Time, m.UserId, m.RoomId, m.MessageType, m.Text)
}

func (m Message) MarshalBinary() ([]byte, error) {
	//TODO: better marshal binary
	return []byte(m.String()), nil
}

func NewMessage(raw string) Message {
	raws := strings.Split(raw, ",")
	parseInt := func(s string) int64 {
		r, _ := strconv.Atoi(s)
		return int64(r)
	}
	msgTime := parseInt(raws[0])
	userId := parseInt(raws[1])
	roomId := raws[2]
	messageType := uint8(parseInt(raws[3]))
	text := raws[4]
	return Message{
		Time:        msgTime,
		UserId:      userId,
		RoomId:      roomId,
		MessageType: messageType,
		Text:        text,
	}

}

func NewNotifyMessage(text string, targetRoom string) *Message {
	msg := Message{
		Time:        time.Now().UnixNano(),
		Text:        text,
		UserId:      0, // system notify
		RoomId:      targetRoom,
		MessageType: MENTAIN_MSG_TYPE,
	}
	return &msg
}
