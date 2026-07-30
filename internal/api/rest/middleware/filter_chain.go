package middleware

import (
	"github.com/Thanga-tamil/noway_service/internal/config"
	"github.com/Thanga-tamil/noway_service/internal/logger"

	"github.com/gin-gonic/gin"
)

func MyMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantId := c.Request.Header.Get("tenant-x")
		path := c.Request.URL.Path
		method := c.Request.Method

		if tenantId == "" {
			logger.Log.Warnw("missing tenant header",
				"path", path,
				"method", method,
				"client_ip", c.ClientIP(),
			)
			c.AbortWithStatusJSON(400, gin.H{"error": "missing tenant-x header"})
			return
		}

		db, ok := config.TenantDBs[tenantId]
		if !ok {
			logger.Log.Warnw("unknown tenant",
				"tenant_id", tenantId,
				"path", path,
				"method", method,
			)
			c.AbortWithStatusJSON(401, gin.H{"error": "invalid tenant"})
			return
		}

		logger.Log.Infow("tenant db resolved",
			"tenant_id", tenantId,
			"path", path,
			"method", method,
		)

		c.Set(tenantId, db)
		c.Next()
	}
}

