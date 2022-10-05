package chat

import (
	"WebSocketServer/app/internal/entity"
	"WebSocketServer/app/internal/repository/chat"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type ServiceI interface {
	Chat(conn *websocket.Conn) error
}

type service struct {
	ChatRepo chat.RepositoryI
}

func NewService(ChatRepo chat.RepositoryI) *service {
	return &service{ChatRepo: ChatRepo}
}

func (s *service) Chat(conn *websocket.Conn) error {
	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		logrus.Error(err)
		return err
	}
	err = s.ChatRepo.NewClient(string(usernameBytes), conn)
	if err != nil {
		logrus.Error(err)
		return err
	}
	logrus.Info("New connection: ", conn.RemoteAddr().String(), " Username: ", string(usernameBytes))

	clients, _ := s.ChatRepo.GetClients()

	for username, conn := range *clients {
		go s.listener(username, conn)
	}

	return nil
}

var pongWait = 5 * time.Second

func (s *service) listener(sender string, conn *websocket.Conn) {
	for {
		//logrus.Info("Listen for: ", conn.RemoteAddr().String())
		//conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		//conn.SetReadDeadline(time.Now().Add(pongWait))
		//conn.SetPongHandler(func(string) error { conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
		_, messageBytes, err := conn.ReadMessage()
		logrus.Info(string(messageBytes))
		message := strings.SplitN(string(messageBytes), ": ", 2) //recipient: text
		if err != nil {
			logrus.Warn("Error during message reading: ", err)
		} else if message[0] != "" && message[1] != "" {
			logrus.Infof("Message from %s to %s: %s", sender, message[0], message[1])
			//writeToAll(message)
			s.writeToOne(message, sender)
		}
	}
}

func (s *service) writeToAll(message string) {
	clients, _ := s.ChatRepo.GetClients()
	for _, conn := range *clients {
		//logrus.Info("Write for: ", conn.RemoteAddr().String())
		err := conn.WriteMessage(1, []byte(message))
		if err != nil {
			logrus.Error("Error during message writing: ", err)
		}
	}
}

func (s *service) writeToOne(message []string, sender string) {
	//logrus.Info("Write for: ", conn.RemoteAddr().String())
	clients, _ := s.ChatRepo.GetClients()
	recipientConn := (*clients)[message[0]]
	if recipientConn == nil || !s.ChatRepo.CheckFriend(&entity.FriendCheck{User: sender, Friend: message[0]}) {
		logrus.Error("Wrong recipient\n")
		return
	}
	err := recipientConn.WriteMessage(1, []byte(sender+": "+message[1]))
	if err != nil {
		logrus.Error("Error during message writing: ", err)
		return
	}
}
