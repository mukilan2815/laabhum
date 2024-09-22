package middleware

import (
	"time"

	"github.com/Mukilan-T/laabhum-gateway-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

func Logger(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		log.Infof("%s | %3d | %13v | %15s | %s %s",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}