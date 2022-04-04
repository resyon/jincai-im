package core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestBackUp_Subscribe(t *testing.T) {
	c := cache.NewRedisClient()
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

func TestBackUp_Subscribe2(t *testing.T) {
	ca := cache.NewRedisClient()
	cb := cache.NewRedisClient()
	cc := cache.NewRedisClient()
	grp := &sync.WaitGroup{}
	grp.Add(2)
	name := "room_name"
	name2 := "room_name_2"

	go func() {
		_, pubSub, err := common.SubUtilReady(ca, name)
		grp.Done()
		assert.NoError(t, err, "fail to sub<sa>")
		err = pubSub.Subscribe(context.TODO(), name2)
		assert.NoError(t, err, "fail to sub<sa> name_2")
		msg := <-pubSub.Channel()
		fmt.Printf("Message from pub_1<sa>: %s\n", msg)
		msg2 := <-pubSub.Channel()
		fmt.Printf("Message from pub_2<sa>: %s\n", msg2)
		grp.Done()
	}()

	go func() {

		_, pubSub, err := common.SubUtilReady(cb, name)
		grp.Done()
		assert.NoError(t, err, "fail to sub<sb>")
		msg := <-pubSub.Channel()
		fmt.Printf("Message from pub_1<sb>: %s\n", msg)
		grp.Done()
	}()

	grp.Wait()
	grp.Add(2)

	cc.Publish(context.TODO(), name2, "test_from_cc_for_name2")
	cc.Publish(context.TODO(), name, "test_from_cc")
}
