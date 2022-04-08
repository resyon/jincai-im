package service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
)

var (
	BackUp = newBakeUp()

	_msgService  = &MessageService{}
	_roomService = &RoomService{}
)

type backUp struct {
	client *redis.Client
	pubSub *redis.PubSub
}

func newBakeUp() *backUp {
	client := cache.NewRedisClient()
	_, pubSub, err := common.SubUtilReady(client, SysChannel)
	if err != nil {
		log.LOG.Panicf("fail to init back up client, Err=%+v", err)
	}
	b := new(backUp)
	b.pubSub = pubSub
	b.client = client

	b.resumeSub()

	go b.backupD()

	return b
}

func (b *backUp) Subscribe(channel string) error {
	return b.pubSub.Subscribe(context.TODO(), channel)
}

func (b *backUp) backupD() {
	for v := range b.pubSub.Channel() {
		log.LOG.Debugf("[IN BACKUP] %#v\n", v)
		msg, _ := model.ParseMessage([]byte(v.Payload))
		err := _msgService.AddMessage(&msg)
		if err != nil {
			log.LOG.Errorf("fail to add message, Err=%+v", err)
		}
	}
}

func (b *backUp) Notify(message *model.Message, channel string) {

	err := b.client.Publish(context.TODO(), channel, message).Err()
	if err != nil {
		log.LOG.Errorf("[Backup] notify: publish, err=%s\n", err)
		return
	}
}

func (b *backUp) resumeSub() {
	room, err := _roomService.GetAllRoom()
	if err != nil {
		log.LOG.Panicf("fail to init backup, Err=%+v", err)
	}
	for _, v := range room {
		err := b.Subscribe(v.RoomId)
		if err != nil {
			log.LOG.Panicf("fail to init backup, Err=%+v", err)
		}
	}
}
