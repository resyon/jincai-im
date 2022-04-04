package model

type Room struct {
	RoomId string `json:"room_id"`
	// id of user who owns the room
	OwnerId  int    `json:"owner_id"`
	RoomName string `json:"room_name"`
}
