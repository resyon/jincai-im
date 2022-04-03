package main

import (
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/middleware"
)

func main() {
	e := gin.Default()
	middleware.EnableAuth(e)
	err := e.Run(":9999")
	if err != nil {
		panic(err)
	}
}
