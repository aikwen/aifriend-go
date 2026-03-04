package mw

import (
	"strconv"
	"time"

	"github.com/aikwen/aifriend-go/pkg/monitor"
	"github.com/gin-gonic/gin"
)

func PrometheusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		path := c.FullPath()
		if path == "" {
			path = "not_found"
		}

		c.Next()

		duration := time.Since(start).Seconds()
		status := strconv.Itoa(c.Writer.Status())
		method := c.Request.Method

		monitor.HttpRequestsTotal.WithLabelValues(method, path, status).Inc()
		monitor.HttpRequestDuration.WithLabelValues(method, path).Observe(duration)
	}
}