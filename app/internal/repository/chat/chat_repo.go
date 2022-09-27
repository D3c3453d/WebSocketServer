package chat

import "github.com/gorilla/websocket"

type RepositoryI interface {
	NewClient(string, *websocket.Conn) error
	GetClients() (*map[string]*websocket.Conn, error)
}

type repository struct {
	clients map[string]*websocket.Conn
}

func NewRepository(clients map[string]*websocket.Conn) *repository {
	return &repository{clients: clients}
}

func (r *repository) NewClient(username string, conn *websocket.Conn) error {
	r.clients[username] = conn
	return nil
}

func (r *repository) GetClients() (*map[string]*websocket.Conn, error) {
	return &r.clients, nil
}
