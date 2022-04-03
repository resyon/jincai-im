package store

import (
	"github.com/resyon/jincai-im/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserDAO_Add(t *testing.T) {
	err := (&UserDAO{}).Add(&model.User{Username: "test", Password: "test"})
	assert.NoError(t, err, "fail to add")
}

func FuzzUserDAO_Add(f *testing.F) {
	dao := &UserDAO{}
	f.Add("test_3", "test")
	f.Fuzz(func(t *testing.T, name, password string) {
		user := &model.User{Username: name, Password: password}
		assert.NoError(t, dao.Add(user), "fail to add user")
		out, err := dao.SelectByUserNameAndPassword(name, password)
		assert.NoError(t, err, "fail to select")
		assert.Equal(t, out.Username, name, "corruption")
		assert.Equal(t, out.Password, password, "corruption")
	})

}
