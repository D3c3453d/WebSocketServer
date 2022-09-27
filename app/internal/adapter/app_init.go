package adapter

import "WebSocketServer/app/internal/adapter/server"

type ServicePool struct {
}

func InitApp() {

	router := server.GetRouter()
	servicePool := InitServices()

	InitHandlers(router, servicePool)

	_ = server.Run(router)
}
