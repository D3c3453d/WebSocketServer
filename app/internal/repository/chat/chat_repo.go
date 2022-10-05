package chat

import (
	"WebSocketServer/app/internal/entity"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type RepositoryI interface {
	NewClient(string, *websocket.Conn) error
	GetClients() (*map[string]*websocket.Conn, error)
	CheckFriend(input *entity.FriendCheck) bool
}

type repository struct {
	clients *map[string]*websocket.Conn
	db      *sqlx.DB
}

func NewRepository(clients *map[string]*websocket.Conn, db *sqlx.DB) *repository {
	return &repository{clients: clients, db: db}
}

func (r *repository) NewClient(username string, conn *websocket.Conn) error {
	(*r.clients)[username] = conn
	return nil
}

func (r *repository) GetClients() (*map[string]*websocket.Conn, error) {
	return r.clients, nil
}

func (r *repository) CheckFriend(input *entity.FriendCheck) bool {
	stmt, err := r.db.PrepareNamed(`SELECT usr.username as user_name, frnd.username as friend_name FROM friends
    									LEFT JOIN users as usr on user_id = usr.id 
    									LEFT JOIN users as frnd on friend_id = frnd.id
    									WHERE usr.username = :user_name AND frnd.username = :friend_name`)
	if err != nil {
		logrus.Error(err)
		return false
	}

	var output entity.FriendCheck
	if err = stmt.Get(&output, &input); err != nil {
		logrus.Error(err)
		return false
	}

	//logrus.Info("output: ", output, "\ninput: ", *input)
	if output == *input {
		return true
	}
	return false
}
