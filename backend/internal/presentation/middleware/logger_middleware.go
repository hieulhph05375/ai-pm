package middleware

import (
	"time"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func JSONLogMiddleware() gin.HandlerFunc {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		entry := logger.WithFields(logrus.Fields{
			"status":     statusCode,
			"latency":    latency.String(),
			"ip":         clientIP,
			"method":     method,
			"path":       path,
			"query":      raw,
			"user_agent": c.Request.UserAgent(),
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			if statusCode >= 400 {
				entry.Warn("request failed")
			} else {
				entry.Info("request handled")
			}
		}
	}
}
