package model

type RoomMate struct {
	RoomId string `json:"room_id" gorm:"room_id"`
	UserId int    `json:"user_id" gorm:"user_id"`
}

func (RoomMate) TableName() string {
	return "room_mate"
}
