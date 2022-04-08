package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/resyon/jincai-im/common"
	"github.com/resyon/jincai-im/log"
	"github.com/resyon/jincai-im/service"
	"net/http"
)

var (
	_msgService = &service.MessageService{}
)

type MessageCtrl struct {
}

func (MessageCtrl) GetHistory(c *gin.Context) {
	id := common.GetUserIdFromContext(c)
	msg, err := _msgService.GetRemoteMessage(id)
	if err != nil {
		info := fmt.Sprintf("fail to get history message, Err=%+v", err)
		log.LOG.Error(info)
		c.JSON(http.StatusInternalServerError, info)
		return
	}
	c.JSON(http.StatusOK, msg)
}
