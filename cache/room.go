package cache

import (
	"context"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/model"
)

const (
	roomMemberKeyPrefix = "room_member_"
	roomCollectionKey   = "room_collection"
)

type RoomCache struct {
}

func (RoomCache) AddRoomToSet(room *model.Room) (err error) {
	err = GetRedis().HSet(context.TODO(), roomCollectionKey, room.RoomId, room.RoomName).Err()
	return
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

func (r *RoomCache) GetAllRoomInfo() ([]*model.Room, error) {
	result, err := GetRedis().HGetAll(context.TODO(), roomCollectionKey).Result()
	if err != nil {
		return nil, err
	}
	ret := make([]*model.Room, len(result))
	i := 0
	for k, v := range result {
		ret[i] = &model.Room{RoomId: k, RoomName: v}
		i++
	}
	return ret, nil
}
