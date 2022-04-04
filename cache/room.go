package cache

import (
	"context"
	"github.com/resyon/jincai-im/common"
)

const (
	roomMemberKeyPrefix   = "room_member_"
	roomCollectionKey     = "room_collection"
	roomNameCollectionKey = "room_name_collection"
)

type RoomCache struct {
}

func (RoomCache) AddRoomToSet(roomName string) (uid string, err error) {
	uid = common.GetUUID()
	err = GetRedis().HSet(context.TODO(), roomNameCollectionKey, roomName, uid).Err()
	err = GetRedis().HSet(context.TODO(), roomCollectionKey, uid, roomName).Err()
	return
}

func (RoomCache) SelectRoomIdByName(roomName string) (uid string, err error) {
	if ok, err := GetRedis().HExists(context.TODO(), roomNameCollectionKey, roomName).Result(); !ok {
		return "", common.RoomNotExistError
	} else if err != nil {
		return "", err
	}
	return GetRedis().HGet(context.TODO(), roomNameCollectionKey, roomName).Result()
}

func (r RoomCache) AddUserToRoom(userId int, roomId string) error {
	if err := r.checkRoomExist(roomId); err != nil {
		return err
	}
	GetRedis().SAdd(context.TODO(), getRoomMemberKey(roomId), userId)
	return nil
}

func (r RoomCache) DelUserFromRoom(userId int, roomId string) error {
	if err := r.checkRoomExist(roomId); err != nil {
		return err
	}
	GetRedis().SRem(context.TODO(), getRoomMemberKey(roomId), userId)
	return nil
}

func (RoomCache) checkRoomExist(roomId string) error {
	exi, err := GetRedis().HExists(context.TODO(), roomCollectionKey, roomId).Result()
	if err != nil {
		return err
	}
	if !exi {
		return common.RoomNotExistError
	}
	return nil
}

func getRoomMemberKey(roomId string) string {
	return roomMemberKeyPrefix + roomId
}
