package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		client := c.ClientIP()
		latency := time.Since(t)
		status := c.Writer.Status()
		path := c.Request.URL.Path
		contentType := c.ContentType()
		staffName := c.GetString("Staffname")
		logrus.WithFields(logrus.Fields{"user": staffName, "client": client, "path": path, "content_type": contentType, "status": status, "latency": latency}).Info("access")
	}
}