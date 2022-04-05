package core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
	"sync"
)

var (
	PeerPool = &peerPool{make(map[int]*PeerConn), &sync.RWMutex{}}
)

const (
	SysChannel = "__sys_channel"
)

type peerPool struct {
	peerMap map[int]*PeerConn
	mutex   *sync.RWMutex
}

func (p *peerPool) AddPeerAndServe(userId int, ws *websocket.Conn) (func(), error) {
	p.mutex.Lock()
	conn, ok := p.peerMap[userId]
	if ok {
		p.mutex.Unlock()
		if err := p.DestroyPeer(userId); err != nil {
			return nil, err
		}
		p.mutex.Lock()
	}
	conn = NewPeerCoon(userId, ws)
	p.peerMap[userId] = conn
	p.mutex.Unlock()

	sub := conn.Client.Subscribe(context.TODO(), SysChannel)
	iFace, err := sub.Receive(context.TODO())
	if err != nil {
		// handle error
		goto Err
	}

	// Should be *Subscription, but others are possible if other actions have been
	// taken on sub since it was created.
	switch v := iFace.(type) {
	case *redis.Subscription:
		// subscribe succeeded
		// ignore
	case *redis.Message:
		// received first message
		conn.ConsumeMessage(v)

	case *redis.Pong:
		// pong received
		// ignore
	default:
		// handle error
		goto Err
	}

	conn.PubSub = sub
	conn.SubPubReady()
	return conn.InitSub(), nil

Err:
	p.mutex.Lock()
	delete(p.peerMap, userId)
	p.mutex.Unlock()
	return nil, err

}

func (p *peerPool) SubscribeChannel(userId int, channel string) error {
	p.mutex.Lock()
	conn, ok := p.peerMap[userId]
	if !ok {
		p.mutex.Unlock()
		return common.RoomNotExistError
	}
	p.mutex.Unlock()
	if conn.PubSub == nil {
		conn.AwaitSubReady()
	}

	return conn.PubSub.Subscribe(conn.ctx, channel)
}

func (p *peerPool) UnSubscribe(userId int, channel string) error {
	p.mutex.Lock()
	conn, ok := p.peerMap[userId]
	if !ok {
		p.mutex.Unlock()
		return common.RoomNotExistError
	}
	p.mutex.Unlock()
	return conn.PubSub.Unsubscribe(conn.ctx, channel)
}

func (p *peerPool) DestroyPeer(userId int) error {
	//TODO: peer collection
	log.LOG.Info("PeerPool#DestroyPeer Called")
	p.mutex.Lock()
	conn, ok := p.peerMap[userId]
	if !ok {
		p.mutex.Unlock()
		return common.RoomNotExistError
	}
	delete(p.peerMap, userId)
	p.mutex.Unlock()
	conn.cancel()
	msg := fmt.Sprintf("%d in room<%s> has downline", conn.userId, SysChannel)
	BackUp.Notify(model.NewNotifyMessage(msg, SysChannel), SysChannel)
	err := conn.PubSub.Close()
	if err != nil {
		return err
	}
	return conn.Client.Close()
}
