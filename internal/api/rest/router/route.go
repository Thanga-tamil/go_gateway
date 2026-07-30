package router

import (
	"github.com/Thanga-tamil/noway_service/internal/api/rest/handler"
	"github.com/Thanga-tamil/noway_service/internal/logger"
	"github.com/gin-gonic/gin"
)

func RouteV1(v1 *gin.RouterGroup)  {

	v1.GET("/ping", func(c *gin.Context) {
		logger.Debug("hello.. ")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	v1.POST("/register", handler.HandleUserRegister)
	v1.GET("/jwt", handler.GenerateJwtToken)

}

func RouteV2(v2 *gin.RouterGroup)  {

	v2.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{ "message": "pong" })
	})

}

