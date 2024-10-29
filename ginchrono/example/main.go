package main

import (
	"log"
	"time"

	"github.com/ashraful88/chrono/ginchrono"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.New()

	// Use default config
	config := ginchrono.DefaultConfig()

	// Customize thresholds if needed
	config.WarningThreshold = 200 * time.Millisecond
	config.ErrorThreshold = 1 * time.Second

	// Enable/disable middleware dynamically
	config.Disable() // Disable logging
	config.Enable()  // Enable logging

	// Custom logger with all requests
	config.LogAllRequests = true
	config.Logger = func(format string, v ...interface{}) {
		log.Printf("[Custom Logger] "+format, v...)
	}

	r.Use(ginchrono.New(config))

	// Example handlers
	r.GET("/fast", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "fast response"})
	})

	r.GET("/slow", func(c *gin.Context) {
		time.Sleep(600 * time.Millisecond) // Will trigger WARNING
		c.JSON(200, gin.H{"message": "slow response"})
	})

	r.GET("/very-slow", func(c *gin.Context) {
		time.Sleep(2100 * time.Millisecond) // Will trigger ERROR
		c.JSON(200, gin.H{"message": "very slow response"})
	})

	r.Run(":8080")
}
