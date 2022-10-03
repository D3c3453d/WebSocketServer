package adapter

import (
	"WebSocketServer/app/internal/adapter/server"
	"WebSocketServer/app/internal/controller"
	chatRepo "WebSocketServer/app/internal/repository/chat"
	"WebSocketServer/app/internal/service/chat"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ServicePool struct {
	ChatService chat.ServiceI
}

func InitApp() {

	router := server.GetRouter()

	clients := make(map[string]*websocket.Conn)
	servicePool := InitServices(clients)

	InitHandlers(router, servicePool)

	_ = server.Run(router)
}

func InitServices(clients map[string]*websocket.Conn) *ServicePool {
	return &ServicePool{
		ChatService: chat.NewService(chatRepo.NewRepository(&clients)),
	}
}

func InitHandlers(router *gin.Engine, pool *ServicePool) {

	publicGroup := router.Group("/")
	privateGroup := router.Group("/")

	userController := controller.NewChatController(pool.ChatService)
	userController.RegisterHandlers(privateGroup, publicGroup)

}
