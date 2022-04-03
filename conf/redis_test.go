package conf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRedisConf(t *testing.T) {
	conf := GetRedisConf()
	assert.NotEmpty(t, conf.Redis.Addr, "fail to get redis config")
}
