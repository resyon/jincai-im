package store

import "github.com/resyon/jincai-im/model"

type RoomMateDAO struct {
}

func (RoomMateDAO) Add(r *model.RoomMate) error {
	return GetDB().Create(r).Error
}

func (RoomMateDAO) SelectRoomIdByUserId(userId int) ([]string, error) {
	var ret []string
	err := GetDB().
		Table(model.RoomMate{}.TableName()).
		Select("room_id").
		Where("user_id=?", userId).
		Find(&ret).Error
	return ret, err
}

func (RoomMateDAO) SelectByRoomId(roomId string) ([]*model.User, error) {
	var ret []*model.User
	err := GetDB().Raw(`
	select distinct(u.id), u.username
	from room_mate r left join user u on u.id = r.user_id
	where room_id=?`, roomId).Find(&ret).Error
	return ret, err
}
