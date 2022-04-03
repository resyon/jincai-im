package store

import (
	"github.com/resyon/jincai-im/model"
)

type UserDAO struct {
}

func (u *UserDAO) Add(user *model.User) error {
	//user.Password = addSalt(user.Password)
	return GetDB().Table(model.User{}.TableName()).Create(user).Error
}

func (u *UserDAO) SelectByUserNameAndPassword(username, password string) (user *model.User, err error) {
	//password = addSalt(password)
	user = new(model.User)
	err = GetDB().Table(model.User{}.TableName()).Where("username=? and password = ?",
		username, password).First(user).Error
	return
}
