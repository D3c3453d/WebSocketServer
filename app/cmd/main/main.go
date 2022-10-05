package main

import (
	"WebSocketServer/app/config"
	"WebSocketServer/app/internal/adapter"
)

func main() {
	config.InitConfig()
	adapter.InitApp()
}
