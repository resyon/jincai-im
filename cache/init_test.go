package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetRedis(t *testing.T) {
	conn := GetRedis()
	assert.NotNil(t, conn, "fail to get a connection")
	resp, err := conn.Ping(context.TODO()).Result()
	assert.NoError(t, err, "fail to ping")
	assert.Equal(t, "PONG", resp, "fail to ping")
}
