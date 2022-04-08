package model

import (
	"strconv"
)

type User struct {
	Id       int    `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Password string `gorm:"column:password"`
}

func (_ User) TableName() string {
	return "user"
}

func (u User) MarshalJSON() ([]byte, error) {
	return []byte(`{"username":"` + u.Username + `","id":` + strconv.Itoa(u.Id) + `}`), nil
}
