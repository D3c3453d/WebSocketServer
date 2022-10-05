package database

import (
	"WebSocketServer/app/config"
	"fmt"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"time"
)

func GetDbConnection() *sqlx.DB {
	dbConfig := config.GetDbConfig()

	cmd := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.Name, dbConfig.SSLMode)

	logrus.Info(dbConfig)

	db, err := sqlx.Open("pgx", cmd)

	if err != nil || db == nil {
		logrus.Panic("Failed init  DB ", err)
		return nil
	}

	start := time.Now()
	for err = db.Ping(); err != nil; {
		logrus.Error("Failed init database ", err)

		if time.Now().After(start.Add(30 * time.Second)) {
			logrus.Panic("Failed init database")
		}

		time.Sleep(5 * time.Second)
	}

	return db
}
