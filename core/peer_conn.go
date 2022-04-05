package core

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
	"io"
)

type PeerConn struct {
	userId    int
	Client    *redis.Client
	PubSub    *redis.PubSub
	pubReady  chan struct{}
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	rcvChan   chan *model.Message
	subTopics []string
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
	// TODO: close handler
	ws.SetCloseHandler(func(code int, text string) error {
		log.LOG.Infof("[WS CLOSE] code=%d, text=%s\n", code, text)
		return PeerPool.DestroyPeer(userId)
	})
	return ret
}

func (p *PeerConn) ConsumeMessage(msg *redis.Message) {
	// TODO: implements
	content := model.NewMessage(msg.Payload)
	log.LOG.Debugf("[PeerConn] read user<%d>: %#v\n", p.userId, content)

	err := p.conn.WriteMessage(websocket.TextMessage, []byte(content.String()))
	if err != nil {
		log.LOG.Debugf("[PeerConn] read user<%d>: %#v\n Err=%#v\n", p.userId, content, err)
	}
}

func (p *PeerConn) PublishMessage(msg *model.Message) {
	log.LOG.Debugf("[PeerConn] send user<%d>: %#v\n", p.userId, msg)
	err := p.Client.Publish(p.ctx, msg.RoomId, msg).Err()
	if err != nil {
		log.LOG.Debugf("[PeerConn] send user<%d>: %#v\n Err=%+v\n", p.userId, msg, err)
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
			log.LOG.Debugf("[PeerConn], type: %d, rcv: %#v\n", msgType, data)

			if err != nil {
				log.LOG.Errorf("Rcv from ws: err=%#v\n", err)
				if err == io.EOF {
					// ignore
					continue
				}

				//TODO: store message when client down
				//if msgType == websocket.CloseMessage {
				err := PeerPool.DestroyPeer(p.userId)
				if err != nil {
					log.LOG.Errorf("fail to destory peer, Err=%+v", err)
				}
				//	return
				//}
				return
			}
			// TODO: consume ping message
			if msgType == websocket.PingMessage {
				continue
			}
			if msgType != websocket.TextMessage {
				log.LOG.Errorf("Rcv from ws: msgType=%d data=%#v\n", msgType, data)
			}

			if msgType == websocket.TextMessage {
				payload := model.NewMessage(string(data))
				payload.UserId = int64(p.userId)
				payload.MessageType = model.COMMON_MSG_TYPE
				p.rcvChan <- &payload
			}
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
				log.LOG.Infof("peer for user<%d> has draw out\n", p.userId)
				return
			}

		}

	}
}
