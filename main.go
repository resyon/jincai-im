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
	room := &controller.RoomCtrl{}

	auth := middleware.EnableAuth(e)

	auth.PATCH("/room", room.JoinRoom)
	auth.POST("/room", room.CreateRoom)
	auth.GET("/ws", room.ServeWS)

	err := e.Run(fmt.Sprintf(":%d", conf.GetAppConf().Port))
	if err != nil {
		log.LOG.Panicf("fail to start serve, Err=%+v", err)
	}
}
