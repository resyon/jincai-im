package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
	"github.com/resyon/jincai-im/store"
)

var (
	_roomCache   = &cache.RoomCache{}
	_roomDao     = &store.RoomDAO{}
	_roomMateDao = &store.RoomMateDAO{}
)

type RoomService struct {
}

func (RoomService) CreateRoom(userId int, roomName string) (*model.Room, error) {
	roomId := common.GetUUID()
	room := &model.Room{
		RoomName: roomName,
		RoomId:   roomId,
		OwnerId:  userId,
	}

	err := _roomDao.AddRoom(room)
	if err != nil {
		return nil, err
	}

	err = _roomCache.AddRoomToSet(room)
	if err != nil {
		return nil, err
	}

	err = BackUp.Subscribe(roomId)
	if err != nil {
		return nil, err
	}

	return room, err
}

func (RoomService) JoinRoom(userId int, roomId string) error {
	err := _roomMateDao.Add(&model.RoomMate{
		RoomId: roomId,
		UserId: userId,
	})
	if err != nil {
		log.LOG.Errorf("fail to insert room_mate, Err=%+v", err)
		return err
	}

	err = _roomCache.AddUserToRoom(userId, roomId)

	// subscribe channel for user
	if err := PeerPool.SubscribeChannel(userId, roomId); err != nil {
		return err
	}

	// notify a user join
	msg := fmt.Sprintf("%d has joined the room %s", userId, roomId)
	log.LOG.Info(msg)
	BackUp.Notify(model.NewNotifyMessage(msg, roomId), roomId)
	return err
}

func (RoomService) GetAllRoom() ([]*model.Room, error) {
	info, err := _roomCache.GetAllRoomInfo()
	if err != nil {
		return nil, err
	}
	if len(info) == 0 {
		room, err := _roomDao.SelectAllRoom()
		if err != nil {
			err = errors.Wrapf(err, "fail to get room info, Err=%+v", err)
			return nil, err
		}
		for _, v := range room {
			info = append(info, v)
			err := _roomCache.AddRoomToSet(v)
			if err != nil {
				err = errors.Wrapf(err, "fail to flush room info, Err=%+v", err)
				log.LOG.Error(err)
			}
		}

	}
	return info, nil
}
