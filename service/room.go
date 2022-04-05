package service

import (
	"fmt"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/core"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
)

var (
	_roomCache = &cache.RoomCache{}
)

type RoomService struct {
}

func (RoomService) CreateRoom(userId int, roomName string) (*model.Room, error) {
	roomId, err := _roomCache.SelectRoomIdByName(roomName)
	if err == nil {
		return &model.Room{RoomId: roomId, RoomName: roomName}, err
	}
	if err != common.RoomNotExistError {
		return nil, err
	}
	roomId, err = _roomCache.AddRoomToSet(roomName)
	if err != nil {
		return nil, err
	}
	err = core.BackUp.Subscribe(roomId)
	if err != nil {
		return nil, err
	}
	room := &model.Room{
		RoomName: roomName,
		RoomId:   roomId,
		OwnerId:  userId,
	}

	return room, err
}

func (RoomService) JoinRoom(userId int, roomId string) error {
	err := _roomCache.AddUserToRoom(userId, roomId)

	// subscribe channel for user
	if err := core.PeerPool.SubscribeChannel(userId, roomId); err != nil {
		return err
	}

	// notify a user join
	msg := fmt.Sprintf("%d has joined the room %s", userId, roomId)
	log.LOG.Info(msg)
	core.BackUp.Notify(model.NewNotifyMessage(msg, roomId), roomId)
	return err
}
