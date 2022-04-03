package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/service"
	"net/http"
)

var (
	_roomService = &service.RoomService{}
)

type RoomCtrl struct {
}

func (RoomCtrl) CreateRoom(c *gin.Context) {
	userId := common.GetUserIdFromContext(c)
	roomName := c.Query("room_name")
	if roomName == "" {
		c.JSON(http.StatusBadRequest, "should bind query named room_name")
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
		}
		return
	}
	//conn, err := (&websocket.Upgrader{
	//	CheckOrigin: func(r *http.Request) bool { return true },
	//}).Upgrade(c.Writer, c.Request, nil)
	//subPub := cache.RedisPool.GetRedisConnection(userId).Subscribe(context.TODO(), roomId)
}
