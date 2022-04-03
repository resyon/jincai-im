package service

import (
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/model"
	"github.com/resyon/jincai-im/store"
)

var (
	_userDAO = &store.UserDAO{}
)

type UserService struct {
}

func (*UserService) Register(username, password string) error {
	password = common.AddSalt(password)
	user := &model.User{Username: username, Password: password}
	return _userDAO.Add(user)
}

func (*UserService) Login(username, password string) (*model.User, error) {
	password = common.AddSalt(password)
	return _userDAO.SelectByUserNameAndPassword(username, password)
}
