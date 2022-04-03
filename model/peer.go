package model

import (
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"sync"
)

var (
	peerOnce  sync.Once
	_peerPool *peerPool
)

type peerPool struct {
	peerMap *sync.Map
	mutex   *sync.RWMutex
}

func (p *peerPool) AddPeer(peer *Peer) {
	p.peerMap.Store(peer.UserId, peer)
}

func GetPeerPool() *peerPool {
	peerOnce.Do(func() {
		_peerPool = &peerPool{
			peerMap: &sync.Map{},
			mutex:   &sync.RWMutex{},
		}
	})
	return _peerPool
}

type Peer struct {
	PubSub *redis.PubSub
	UserId int
	Conn   *websocket.Conn
}
