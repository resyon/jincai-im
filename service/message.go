package service

import (
	"github.com/resyon/jincai-im/model"
	"github.com/resyon/jincai-im/store"
)

var (
	_msgDao = &store.MessageDAO{}
)

type MessageService struct {
}

func (MessageService) AddMessage(msg *model.Message) error {
	return _msgDao.AddMessage(msg)
}

func (MessageService) GetRemoteMessage(userId int) ([]*model.Message, error) {
	return _msgDao.SelectMessageByUserId(userId)
}
