package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func wsHandler(ctx *gin.Context) {
	conn, err := wsUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logrus.Errorln("Error during connection upgrade: ", err)
		return
	}
	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logrus.Errorln("Connection close error: ", err)
			return
		}
	}(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Errorln("Error during message reading: ", err)
			break
		}
		logrus.Infof("Received: %s", message)
		err = conn.WriteMessage(messageType, message)
		if err != nil {
			logrus.Errorln("Error during message writing: ", err)
			break
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ws", wsHandler)

	err := r.Run("localhost:7077")
	if err != nil {
		logrus.Errorln("Error run server: ", err)
		return
	}
}
