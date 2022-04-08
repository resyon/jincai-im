package service

import (
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
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
	conn = NewPeerCoonNoWait(userId, ws)
	p.peerMap[userId] = conn
	p.mutex.Unlock()

	err := conn.ReadyForSubscribe()

	if err != nil {
		goto Err
	}
	return conn.InitAndGetServer(), nil

Err:
	p.mutex.Lock()
	delete(p.peerMap, userId)
	p.mutex.Unlock()
	return nil, err

}

func (p *peerPool) SubscribeChannel(userId int, channel string) error {
	p.mutex.RLock()
	conn, ok := p.peerMap[userId]
	if !ok {
		p.mutex.RUnlock()
		return common.RoomNotExistError
	}
	p.mutex.RUnlock()

	return conn.Subscribe(conn.ctx, channel)
}

func (p *peerPool) UnSubscribe(userId int, channel string) error {
	p.mutex.RLock()
	conn, ok := p.peerMap[userId]
	if !ok {
		p.mutex.RUnlock()
		return common.RoomNotExistError
	}
	p.mutex.RUnlock()
	return conn.Unsubscribe(conn.ctx, channel)
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
	return conn.Destroy()
}
