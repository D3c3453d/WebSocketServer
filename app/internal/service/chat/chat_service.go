package chat

import (
	"WebSocketServer/app/internal/repository/chat"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"strings"
)

type ServiceI interface {
	Chat(conn *websocket.Conn) error
}

type service struct {
	ChatRepo chat.RepositoryI
}

func NewService() *service {
	return &service{}
}

func (s *service) Chat(conn *websocket.Conn) error {
	_, usernameBytes, err := conn.ReadMessage()
	if err != nil {
		logrus.Error(err)
		return err
	}
	s.ChatRepo.NewClient(string(usernameBytes), conn)
	logrus.Info("New connection: ", conn.RemoteAddr().String(), " Username: ", string(usernameBytes))

	clients, _ := s.ChatRepo.GetClients()

	for username, conn := range *clients {
		go s.listener(username, conn)
	}

	return nil
}

func (s *service) listener(sender string, conn *websocket.Conn) {
	for {
		//logrus.Info("Listen for: ", conn.RemoteAddr().String())
		//conn.SetReadDeadline(time.Now().Add(time.Second * 1))
		_, messageBytes, err := conn.ReadMessage()
		message := strings.SplitN(string(messageBytes), ": ", 2) //recipient: text
		if err != nil {
			logrus.Warn("Error during message reading: ", err)
		} else if message[0] != "" && message[1] != "" {
			logrus.Infof("Received from %s: %s", conn.RemoteAddr().String(), message)
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
	if recipientConn == nil {
		logrus.Error("Wrong recipient: ")
		return
	}
	err := recipientConn.WriteMessage(1, []byte(sender+": "+message[1]))
	if err != nil {
		logrus.Error("Error during message writing: ", err)
		return
	}
}
