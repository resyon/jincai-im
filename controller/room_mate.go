package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/log"
	"net/http"
)

type RoomMateCtrl struct {
}

func (RoomMateCtrl) GetRoomMate(c *gin.Context) {
	roomId := c.Query("room_id")
	if roomId == "" {
		c.JSON(http.StatusBadRequest, "Empty room_id")
		return
	}
	mate, err := _roomMateService.GetRoomMate(roomId)
	if err != nil {
		info := fmt.Sprintf("fail to get room-mate, Err=%+v", err)
		log.LOG.Error(info)
		c.JSON(http.StatusInternalServerError, info)
		return
	}
	c.JSON(http.StatusOK, mate)
}
