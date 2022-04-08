package service

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/cache"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/model"
	"io"
)

type PeerConn struct {
	userId    int
	client    *redis.Client
	pubSub    *redis.PubSub
	pubReady  chan struct{}
	conn      *websocket.Conn
	ctx       context.Context
	cancel    context.CancelFunc
	rcvChan   chan *model.Message
	subTopics []string
}

func NewPeerCoon(userId int, ws *websocket.Conn) (*PeerConn, error) {
	conn := NewPeerCoonNoWait(userId, ws)
	if err := conn.ReadyForSubscribe(); err != nil {
		log.LOG.Errorf("Fail to subscribe, Err=%+v", err)
		return nil, err
	}
	return conn, nil
}

func (conn *PeerConn) Destroy() error {

	conn.cancel()
	msg := fmt.Sprintf("%d in room<%s> has downline", conn.userId, SysChannel)
	BackUp.Notify(model.NewNotifyMessage(msg, SysChannel), SysChannel)

	err := conn.pubSub.Close()
	if err != nil {
		return err
	}
	return conn.client.Close()
}

// NewPeerCoonNoWait 初始化连接, 立即返回, 需手动订阅系统信道
func NewPeerCoonNoWait(userId int, ws *websocket.Conn) *PeerConn {
	ctx, cancel := context.WithCancel(context.TODO())
	ret := &PeerConn{
		ctx:      ctx,
		cancel:   cancel,
		conn:     ws,
		userId:   userId,
		rcvChan:  make(chan *model.Message),
		pubReady: make(chan struct{}, 1),
		client:   cache.NewRedisClient(),
	}
	// TODO: close handler
	ws.SetCloseHandler(func(code int, text string) error {
		log.LOG.Infof("[WS CLOSE] code=%d, text=%s\n", code, text)
		return PeerPool.DestroyPeer(userId)
	})
	return ret
}

func (conn *PeerConn) ReadyForSubscribe() error {
	sub := conn.client.Subscribe(context.TODO(), SysChannel)
	iFace, err := sub.Receive(context.TODO())
	if err != nil {
		// handle error
		return err
	}

	// Should be *Subscription, but others are possible if other actions have been
	// taken on sub since it was created.
	switch v := iFace.(type) {
	case *redis.Subscription:
		// subscribe succeeded
		// ignore
	case *redis.Message:
		// received first message
		conn.NotifyClient(v)

	case *redis.Pong:
		// pong received
		// ignore
	default:
		// handle error
		return err
	}

	conn.pubSub = sub
	conn.SubPubReady()

	return nil
}

func (conn *PeerConn) Subscribe(ctx context.Context, channel string) error {
	if conn.pubSub == nil {
		<-conn.pubReady
	}
	return conn.pubSub.Subscribe(ctx, channel)
}

func (conn *PeerConn) Unsubscribe(ctx context.Context, channel string) error {
	if conn.pubSub == nil {
		<-conn.pubReady
	}
	return conn.pubSub.Unsubscribe(ctx, channel)
}

func (conn *PeerConn) BoardCast(msg *model.Message) {
	// TODO: implements

	msg.UserId = conn.userId
	msg.Id = common.GenerateID()
	msg.MessageType = model.CommonMsgType
	conn.PublishMessage(msg)
}

func (conn *PeerConn) NotifyClient(msg *redis.Message) {

	content, err := model.ParseMessage([]byte(msg.Payload))
	if err != nil {
		log.LOG.Infof("broken message, Err=%+v", err)
		return
	}
	log.LOG.Debugf("[PeerConn] read user<%d>: %#v\n", conn.userId, content)
	err = conn.conn.WriteMessage(websocket.TextMessage, []byte(msg.String()))
	if err != nil {
		log.LOG.Errorf("[PeerConn] notify user<%d>: %#v\n Err=%#v\n", conn.userId, msg, err)
	}
}

func (conn *PeerConn) PublishMessage(msg *model.Message) {
	log.LOG.Debugf("[PeerConn] publish user<%d>: %#v\n", conn.userId, msg)
	err := conn.client.Publish(conn.ctx, msg.RoomId, msg).Err()
	if err != nil {
		log.LOG.Debugf("[PeerConn] publish user<%d>: %#v\n Err=%+v\n", conn.userId, msg, err)
	}
}

func (conn *PeerConn) SubPubReady() {
	conn.pubReady <- struct{}{}
}

func (conn *PeerConn) AwaitSubReady() {
	<-conn.pubReady
}

func (conn *PeerConn) InitAndGetServer() (serveFunc func()) {
	ch := conn.pubSub.Channel()

	go func() {
		for {
			msgType, data, err := conn.conn.ReadMessage()
			log.LOG.Debugf("[PeerConn], type: %d, rcv: %#v\n", msgType, data)

			if err != nil {
				log.LOG.Errorf("Rcv from ws: err=%#v\n", err)
				if err == io.EOF {
					// ignore
					continue
				}

				//TODO: store message when client down
				//if msgType == websocket.CloseMessage {
				err := PeerPool.DestroyPeer(conn.userId)
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
				payload, err := model.ParseMessage(data)
				if err != nil {
					continue
				}
				payload.UserId = conn.userId
				payload.MessageType = model.CommonMsgType
				conn.rcvChan <- &payload
			}
		}
	}()

	return func() {

		for {

			select {

			// read
			case msg := <-ch:
				conn.NotifyClient(msg)

			// write
			case msg := <-conn.rcvChan:
				conn.BoardCast(msg)

			// cancel
			case <-conn.ctx.Done():
				log.LOG.Infof("peer for user<%d> has draw out\n", conn.userId)
				return
			}

		}

	}
}
