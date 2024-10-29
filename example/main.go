package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/ginchrono"
)

func main() {
	r := gin.New()

	// Use default config
	config := ginchrono.DefaultConfig()
	r.Use(ginchrono.New(config))

	// Or use custom config
	customConfig := ginchrono.Config{
		WarningThreshold: 200 * time.Millisecond,
		LogAllRequests:  true,
		Logger: func(format string, v ...interface{}) {
			log.Printf("[Custom Logger] "+format, v...)
		},
	}
	r.Use(ginchrono.New(customConfig))

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	r.Run(":8080")
}
