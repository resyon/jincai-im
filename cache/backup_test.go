package cache

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/common"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackUp_Subscribe(t *testing.T) {
	c := NewRedisClient()
	cond := make(chan struct{})
	s1 := "t1"
	s2 := "t2"
	var pubSub *redis.PubSub
	go func() {
		var err error
		_, pubSub, err = common.SubUtilReady(c, s1)
		cond <- struct{}{}
		assert.NoError(t, err, "fail to sub")
		for msg := range pubSub.Channel() {
			fmt.Printf("Message from pub_1<s1>: %s\n", msg)
		}
	}()

	<-cond
	c.Publish(context.TODO(), s1, "t1")
	err := pubSub.Subscribe(context.TODO(), s2)
	assert.NoError(t, err, "fail to sub_2")
	c.Publish(context.TODO(), s2, "t2")
	c.Publish(context.TODO(), s1, "t3")
}
