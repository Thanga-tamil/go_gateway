package app

import (
	"github.com/Thanga-tamil/noway_service/internal/config"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/Thanga-tamil/noway_service/internal/service"
)

func App(c config.Cfg) {

	logger.Info("Initialize required services from app.go")

	config.InitSql(c)

	if pong, err := config.InitRedis(c); err != nil {
		logger.Fatalf("Error connecting to Redis: %s", err)
	} else {
		logger.Infof("Connected to Redis: %s Redis init success", pong)
	}

	if err := service.LoadJwtSignKeyInCache(); err != nil {
		logger.Fatalf("Error while loading Jwt sign key in inmemory: %s", err)
	} else {
		logger.Info("Jwt sign key loaded in memory successfully")
	}

	config.InitRabbitMQConnection()
 
	logger.Info("Required services initialization completed successfully")
}

