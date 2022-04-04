package core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/model"
	"log"
)

var (
	BackUp = newBakeUp()
)

type backUp struct {
	client *redis.Client
	pubSub *redis.PubSub
}

func newBakeUp() *backUp {
	client := cache.NewRedisClient()
	_, pubSub, err := common.SubUtilReady(client, SysChannel)
	if err != nil {
		panic(err)
	}
	b := new(backUp)
	b.pubSub = pubSub
	b.client = client

	go b.backupD()

	return b
}

func (b *backUp) Subscribe(channel string) error {
	return b.pubSub.Subscribe(context.TODO(), channel)
}

func (b *backUp) backupD() {
	for v := range b.pubSub.Channel() {
		//TODO: persist message
		log.Printf("[IN BACKUP] %#v\n", v)
	}
}

func (b *backUp) Notify(message *model.Message, channel string) {

	//TODO: handle the error
	err := b.client.Publish(context.TODO(), channel, message).Err()
	if err != nil {
		fmt.Printf("[Backup] notify: publish, err=%s\n", err)
		return
	}
}
