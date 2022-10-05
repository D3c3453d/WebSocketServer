package server

import (
	"WebSocketServer/app/config"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(router *gin.Engine) error {

	serverConfig := config.GetServerConfig()

	server := &http.Server{
		Addr:    serverConfig.Bind + ":" + serverConfig.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: " + err.Error())
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()
	return server.Shutdown(ctx)
}

func GetRouter() *gin.Engine {
	var router *gin.Engine

	router = gin.Default()
	//router.Use(enableCors())

	return router
}
