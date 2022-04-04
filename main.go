package main

import (
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/controller"
	"github.com/resyon/jincai-im/middleware"
)

func main() {
	e := gin.Default()
	room := &controller.RoomCtrl{}

	auth := middleware.EnableAuth(e)

	auth.PATCH("/room", room.JoinRoom)
	auth.POST("/room", room.CreateRoom)
	auth.GET("/ws", room.ServeWS)

	err := e.Run(":9999")
	if err != nil {
		panic(err)
	}
}
