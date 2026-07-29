package app

import (
	"github.com/sirupsen/logrus"

	"github.com/Thanga-tamil/noway_service/internal/config"
	"github.com/Thanga-tamil/noway_service/internal/service"
)

func App(c config.Cfg) {

	logrus.Info("Initialize required services from app.go")

	// config.InitSql(c)

	// if pong, err := config.InitRedis(c); err != nil {
	// 	logrus.Fatalf("Error connecting to Redis: %s", err)
	// } else {
	// 	logrus.Infof("Connected to Redis: %s Redis init success", pong)
	// }

	if err := service.LoadJwtSignKeyInCache(); err != nil {
		logrus.Fatalf("Error while loading Jwt sign key in inmemory: %s", err)
	} else {
		logrus.Info("Jwt sign key loaded in memory successfully")
	}
 
	logrus.Info("Required services initialization completed successfully")
}
