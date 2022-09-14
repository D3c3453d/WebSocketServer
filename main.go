package main

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"net/http"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var clients = make(map[*websocket.Conn]struct{})

//var clients = make(map[string]*websocket.Conn)

func writeForAll(messageType int, message []byte) {
	for conn := range clients {
		logrus.Info("Write for: ", conn.RemoteAddr().String())
		err := conn.WriteMessage(messageType, message)
		if err != nil {
			logrus.Error("Error during message writing: ", err)
		}
	}
}

func listener(conn *websocket.Conn) {
	for {
		logrus.Info("Listen for: ", conn.RemoteAddr().String())
		//conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Warn("Error during message reading: ", err)
		} else if message != nil {
			logrus.Infof("Received from %s: %s", conn.RemoteAddr().String(), message)
			writeForAll(messageType, message)
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Error during connection upgrade: ", err)
		return
	} else {
		logrus.Info("New connection: ", conn.RemoteAddr().String())
	}

	//clients[strconv.Itoa(int(time.Now().Unix()))] = conn
	clients[conn] = struct{}{}
	for conn := range clients {
		go listener(conn)
	}
}

func main() {
	http.HandleFunc("/ws", wsHandler)
	logrus.Fatal(http.ListenAndServe("localhost:7077", nil))
	//r := gin.Default()
	//
	//r.GET("/ws", wsHandler)
	//
	//err := r.Run("localhost:7077")
	//if err != nil {
	//	logrus.Error("Error run server: ", err)
	//	return
	//}
}
