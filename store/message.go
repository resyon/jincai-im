package store

import (
	"github.com/resyon/jincai-im/model"
)

type MessageDAO struct {
}

func (MessageDAO) AddMessage(msg *model.Message) error {
	return GetDB().Table(msg.TableName()).Create(msg).Error
}

func (MessageDAO) SelectMessageByRoomId(roomId string) ([]*model.Message, error) {
	//TODO: add page split
	var ret []*model.Message

	err := GetDB().
		Table(model.Message{}.TableName()).
		Where("room_id=?", roomId).
		Order("time").
		Find(&ret).Error
	return ret, err
}

func (MessageDAO) SelectMessageByUserId(userId int) ([]*model.Message, error) {
	//TODO: add page split
	var ret []*model.Message

	err := GetDB().
		Table(model.Message{}.TableName()).
		Where("user_id=?", userId).
		Order("time").
		Find(&ret).Error
	return ret, err
}
