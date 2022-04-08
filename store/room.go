package store

import "github.com/resyon/jincai-im/model"

type RoomDAO struct {
}

func (RoomDAO) AddRoom(room *model.Room) error {
	return GetDB().Create(room).Error
}

func (RoomDAO) SelectRoomById(id string) (*model.Room, error) {
	ret := new(model.Room)
	err := GetDB().Where("room_id=?", id).First(ret).Error
	return ret, err
}

func (RoomDAO) SelectAllRoom() ([]*model.Room, error) {
	var ret []*model.Room
	err := GetDB().Table(model.Room{}.TableName()).Find(&ret).Error
	return ret, err
}
