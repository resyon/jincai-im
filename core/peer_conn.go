package core

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/model"
	"io"
)

type PeerConn struct {
	userId   int
	Client   *redis.Client
	PubSub   *redis.PubSub
	pubReady chan struct{}
	conn     *websocket.Conn
	ctx      context.Context
	cancel   context.CancelFunc
	rcvChan  chan *model.Message
}

func NewPeerCoon(userId int, ws *websocket.Conn) *PeerConn {
	ctx, cancel := context.WithCancel(context.TODO())
	ret := &PeerConn{
		ctx:      ctx,
		cancel:   cancel,
		conn:     ws,
		userId:   userId,
		rcvChan:  make(chan *model.Message),
		pubReady: make(chan struct{}, 1),
		Client:   cache.NewRedisClient(),
	}
	ws.SetCloseHandler(func(code int, text string) error {
		fmt.Printf("[WS CLOSE] code=%d, text=%s\n", code, text)
		return PeerPool.DestroyPeer(userId)
	})
	return ret
}

func (p *PeerConn) ConsumeMessage(msg *redis.Message) {
	// TODO: implements
	content := model.NewMessage(msg.Payload)
	fmt.Printf("[PeerConn] read user<%d>: %#v\n", p.userId, content)

	err := p.conn.WriteMessage(websocket.TextMessage, []byte(content.String()))
	if err != nil {
		fmt.Printf("[PeerConn] read user<%d>: %#v\n Err=%#v\n", p.userId, content, err)
	}
}

func (p *PeerConn) PublishMessage(msg *model.Message) {
	fmt.Printf("[PeerConn] send user<%d>: %#v\n", p.userId, msg)
	err := p.Client.Publish(p.ctx, msg.RoomId, msg).Err()
	if err != nil {
		fmt.Printf("[PeerConn] send user<%d>: %#v\n Err=%+v\n", p.userId, msg, err)
	}
}

func (p *PeerConn) SubPubReady() {
	p.pubReady <- struct{}{}
}

func (p *PeerConn) AwaitSubReady() {
	<-p.pubReady
}

func (p *PeerConn) InitSub() (serveFunc func()) {
	ch := p.PubSub.Channel()

	go func() {
		for {
			msgType, data, err := p.conn.ReadMessage()
			//fmt.Printf("[PeerConn], type: %d, rcv: %#v\n", msgType, data)
			if err != nil {
				if err == io.EOF {
					// ignore
					continue
				}
				fmt.Printf("Rcv from ws: err=%#v\n", err)
				return
			}
			if msgType != websocket.TextMessage {
				fmt.Printf("Rcv from ws: msgType=%d data=%#v\n", msgType, data)
				continue
			}

			payload := model.NewMessage(string(data))
			payload.MessageType = model.COMMON_MSG_TYPE
			p.rcvChan <- &payload
		}
	}()

	return func() {

		for {

			select {

			// read
			case msg := <-ch:
				p.ConsumeMessage(msg)

			// write
			case msg := <-p.rcvChan:
				p.PublishMessage(msg)

			// cancel
			case <-p.ctx.Done():
				fmt.Printf("peer for user<%d> has draw out\n", p.userId)
				return
			}

		}

	}
}
