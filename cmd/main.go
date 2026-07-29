package main

import (
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"github.com/Thanga-tamil/noway_service/internal/api/rest/middleware"
	"github.com/Thanga-tamil/noway_service/internal/api/rest/router"
	"github.com/Thanga-tamil/noway_service/internal/app"
	"github.com/Thanga-tamil/noway_service/internal/config"
)


var (
    WarningLog *log.Logger
    InfoLog   *log.Logger
    ErrorLog   *log.Logger
)

func main() {

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		DisableColors:   false,
		// TimestampFormat: "2006-01-02 15:04:05",
	})


	// func OpenFile(name string, flag int, perm FileMode) (*File, error) {
	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("Error while opening service log append file: ", err.Error())
	}
	defer logFile.Close()

    logrus.SetOutput(logFile)

	logrus.Info("Serve http request response using engine 'net/http'")

	conf := config.LoadConfig()

	app.App(conf)
	
	r := gin.Default() // r := gin.New()

	v1 := r.Group("/api/v1")
	v2 := r.Group("/api/v2")

	v1.Use(middleware.MyMiddleware())
	
	router.RouteV1(v1)
	router.RouteV2(v2)

	r.Run(conf.Host + ":" + strconv.Itoa(conf.Port))

}


