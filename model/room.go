package model

type Room struct {
	RoomId string `json:"room_id" gorm:"room_id"`
	// id of user who owns the room
	OwnerId  int    `json:"owner_id" gorm:"owner_id"`
	RoomName string `json:"room_name" gorm:"room_name"`
}

func (Room) TableName() string {
	return "room"
}
