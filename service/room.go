package service

import (
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/model"
)

var (
	_roomCache = &cache.RoomCache{}
)

type RoomService struct {
}

func (RoomService) CreateRoom(userId int, roomName string) (*model.Room, error) {
	roomId, err := _roomCache.AddRoomToSet(roomName)
	if err != nil {
		return nil, err
	}
	err = cache.BackUp.Subscribe(roomId)
	if err != nil {
		return nil, err
	}
	room := &model.Room{
		RoomName: roomName,
		RoomId:   roomId,
		OwnerId:  userId,
	}
	model.GetRoomPool().AddRoom(room)

	return room, err
}

func (RoomService) JoinRoom(userId int, roomId string) error {
	err := _roomCache.AddUserToRoom(userId, roomId)
	
	return err
}
