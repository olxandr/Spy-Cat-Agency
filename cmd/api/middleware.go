package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		end := time.Now()
		latency := end.Sub(start)

		// Log request details
		log.Printf("Method: %s | Status: %d | Latency: %s | ClientIP: %s | Path: %s\n",
			c.Request.Method,
			c.Writer.Status(),
			latency,
			c.ClientIP(),
			c.Request.URL.Path,
		)
	}
}
