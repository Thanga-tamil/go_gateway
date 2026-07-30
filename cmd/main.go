package main

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/Thanga-tamil/noway_service/internal/api/rest/middleware"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/Thanga-tamil/noway_service/internal/api/rest/router"
	"github.com/Thanga-tamil/noway_service/internal/app"
	"github.com/Thanga-tamil/noway_service/internal/config"
)


func main() {

	conf := config.LoadConfig()

	logger.Init("app.log", conf.LogLevel)

	logger.Infof("Starting HTTP server using net/http engine")

	app.App(conf)
	
	r := gin.Default()

	v1 := r.Group("/api/v1")
	v2 := r.Group("/api/v2")

	v1.Use(middleware.MyMiddleware())
	
	router.RouteV1(v1)
	router.RouteV2(v2)

	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))

}



