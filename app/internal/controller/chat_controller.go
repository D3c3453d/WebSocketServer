package controller

import (
	"WebSocketServer/app/internal/service/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type ChatController struct {
	ChatService chat.ServiceI
}

func NewChatController(ChatService chat.ServiceI) *ChatController {
	return &ChatController{ChatService: ChatService}
}

func (c *ChatController) RegisterHandlers(_, public *gin.RouterGroup) {
	public.GET("/ws", c.wsHandler)
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func (c *ChatController) wsHandler(ctx *gin.Context) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Error("Error during connection upgrade: ", err)
		return
	}
	err = c.ChatService.Chat(conn)
	if err != nil {
		logrus.Error(err)
		return
	}
}
