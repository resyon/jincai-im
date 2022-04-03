package model

const (
	MENTAIN_MSG_TYPE = iota
	COMMON_MSG_TYPE
	HEARTBEAT_MSG_TYPE
)

type Message struct {
	Time        int64  `json:"time"`
	Text        string `json:"text"`
	UserId      int64  `json:"user_id"`
	RoomId      string `json:"room_id"` // 64B
	MessageType uint8  `json:"message_type"`
}

func (_ Message) TableName() string {
	return "message"
}
