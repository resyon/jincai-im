package service

import "github.com/resyon/jincai-im/model"

type RoomMateService struct {
}

func (RoomMateService) GetJoinedRoom(userId int) ([]string, error) {
	return _roomMateDao.SelectRoomIdByUserId(userId)
}

func (RoomMateService) GetRoomMate(roomId string) ([]*model.User, error) {
	return _roomMateDao.SelectByRoomId(roomId)
}
