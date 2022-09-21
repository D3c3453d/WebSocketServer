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

//var clients = make(map[*websocket.Conn]struct{})

var clients = make(map[string]*websocket.Conn)

func writeToAll(message []byte) {
	for _, conn := range clients {
		//logrus.Info("Write for: ", conn.RemoteAddr().String())
		err := conn.WriteMessage(1, message)
		if err != nil {
			logrus.Error("Error during message writing: ", err)
		}
	}
}

func writeToOne(message []byte, sender string, recipient string) {
	//logrus.Info("Write for: ", conn.RemoteAddr().String())
	recipientConn := clients[recipient]
	if recipientConn == nil {
		logrus.Error("Wrong recipient: ")
		return
	}
	err := recipientConn.WriteMessage(1, []byte(sender))
	err = recipientConn.WriteMessage(1, message)
	if err != nil {
		logrus.Error("Error during message writing: ", err)
		return
	}
}

func listener(sender string, conn *websocket.Conn) {
	for {
		//logrus.Info("Listen for: ", conn.RemoteAddr().String())
		//conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		_, username, err := conn.ReadMessage()
		_, message, err := conn.ReadMessage()
		if err != nil {
			logrus.Warn("Error during message reading: ", err)
		} else if username != nil && message != nil {
			logrus.Infof("Received from %s: %s", conn.RemoteAddr().String(), message)
			//writeToAll(message)
			writeToOne(message, sender, string(username))
		}
	}
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Error("Error during connection upgrade: ", err)
		return
	}

	_, usernameBytes, err := conn.ReadMessage()
	logrus.Info("New connection: ", conn.RemoteAddr().String(), " Username: ", string(usernameBytes))
	clients[string(usernameBytes)] = conn
	for username, conn := range clients {
		go listener(username, conn)
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
