package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func Run(router *gin.Engine) error {

	server := &http.Server{
		Addr:    "localhost:7077",
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			logrus.Fatalf("Failed to listen and serve: " + err.Error())
		}
	}()

	quet := make(chan os.Signal, 1)
	signal.Notify(quet, os.Interrupt, os.Interrupt)

	<-quet

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
