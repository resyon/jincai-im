package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/conf"
	"github.com/resyon/jincai-im/controller"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/middleware"
)

func main() {
	e := gin.Default()

	auth := middleware.EnableAuth(e)

	room := &controller.RoomCtrl{}
	auth.PATCH("/room", room.JoinRoom)
	auth.POST("/room", room.CreateRoom)
	auth.GET("/room", room.GetAllRoom)
	auth.GET("/ws", room.ServeWS)

	msg := &controller.MessageCtrl{}
	auth.GET("/message", msg.GetHistory)

	roomMate := &controller.RoomMateCtrl{}
	auth.GET("/room_mate", roomMate.GetRoomMate)

	err := e.Run(fmt.Sprintf(":%d", conf.GetAppConf().Port))
	if err != nil {
		log.LOG.Panicf("fail to start serve, Err=%+v", err)
	}
}
