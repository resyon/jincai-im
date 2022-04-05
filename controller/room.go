package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/core"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/service"
	"net/http"
)

var (
	_roomService = &service.RoomService{}
	upgrader     = &websocket.Upgrader{
		CheckOrigin: func(_ *http.Request) bool {
			return true
		},
	}
)

type RoomCtrl struct {
}

func (RoomCtrl) CreateRoom(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	roomName := c.Query("room_name")
	if roomName == "" {
		log.LOG.Info("empty room name")
		c.JSON(http.StatusBadRequest, "should bind query named room_name")
		return
	}
	room, err := _roomService.CreateRoom(userId, roomName)
	if err != nil {
		log.LOG.Errorf("fail to create room, Err=%+v\n", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, room)
}

func (RoomCtrl) JoinRoom(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	roomId := c.Query("room_id")
	if roomId == "" {
		log.LOG.Info("empty room name")
		c.JSON(http.StatusBadRequest, "should bind query named room_id")
		return
	}
	err := _roomService.JoinRoom(userId, roomId)
	if err != nil {
		if err == common.RoomNotExistError {
			log.LOG.Info("join a not exist room")
			c.JSON(http.StatusBadRequest, "invalid room_id, may be you need create a room")
			return
		}
		log.LOG.Errorf("fail to join room, Err=%+v\n", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userId})
}

func (RoomCtrl) ServeWS(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	// TODO: fallback to http
	if err != nil {
		log.LOG.Errorf("fail to upgrade to websocket, Err=%+v\n", err)
		c.JSON(http.StatusNotImplemented, "websocket unsupported "+err.Error())
		return
	}

	server, err := core.PeerPool.AddPeerAndServe(userId, conn)
	if err != nil {
		log.LOG.Errorf("fail to add peer, Err=%+v\n", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	server()
}
