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
		logrus.Error("Error during connection upgrade: ", err)
		return
	}

	var clients map[*websocket.Conn]bool
	clients[conn] = true

	defer func(conn *websocket.Conn) {
		err := conn.Close()
		if err != nil {
			logrus.Error("Connection close error: ", err)
			return
		}
		delete(clients, conn)
	}(conn)

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Error("Error during message reading: ", err)
		}
		logrus.Infof("Received: %s", message)
		for conn := range clients {
			err := conn.WriteMessage(messageType, message)
			if err != nil {
				logrus.Error("Error during message writing: ", err)
			}
		}
		if err != nil {
			logrus.Error("Error during message writing: ", err)
		}
	}
}

func main() {
	r := gin.Default()

	r.GET("/ws", wsHandler)

	err := r.Run("localhost:7077")
	if err != nil {
		logrus.Error("Error run server: ", err)
		return
	}
}
