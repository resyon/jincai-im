package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/core"
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
		c.JSON(http.StatusBadRequest, "should bind query named room_name")
		return
	}
	room, err := _roomService.CreateRoom(userId, roomName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, room)
}

func (RoomCtrl) JoinRoom(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	roomId := c.Query("room_id")
	if roomId == "" {
		c.JSON(http.StatusBadRequest, "should bind query named room_id")
		return
	}
	err := _roomService.JoinRoom(userId, roomId)
	if err != nil {
		if err == common.RoomNotExistError {
			c.JSON(http.StatusBadRequest, "invalid room_id, may be you need create a room")
			return
		}
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
		c.JSON(http.StatusNotImplemented, "websocket unsupported "+err.Error())
		return
	}

	server, err := core.PeerPool.AddPeerAndServe(userId, conn)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	server()
}
