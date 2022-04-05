package model

type User struct {
	Id       int    `json:"id" gorm:"column:id"`
	Username string `json:"username" gorm:"column:username"`
	Password string `json:"password" gorm:"column:password"`
}

func (_ User) TableName() string {
	return "user"
}
