package ginchrono

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// LogLevel represents different logging levels
type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	ERROR LogLevel = "ERROR"
)

// Config holds configuration for the middleware
type Config struct {
	// Enable/disable the middleware
	Enabled bool
	// Threshold in milliseconds after which to log warning
	WarningThreshold time.Duration
	// Threshold in milliseconds after which to log error
	ErrorThreshold time.Duration
	// Whether to log all requests, not just slow ones
	LogAllRequests bool
	// Custom logger function - if nil, standard log package is used
	Logger func(format string, v ...interface{})
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Enabled:          true,
		WarningThreshold: 500 * time.Millisecond,  // 500ms default for warnings
		ErrorThreshold:   2000 * time.Millisecond, // 2s default for errors
		LogAllRequests:   false,                   // only log slow requests by default
		Logger:           log.Printf,              // use standard logger by default
	}
}

// determineLogLevel determines the log level based on duration and thresholds
func determineLogLevel(duration time.Duration, config Config) LogLevel {
	switch {
	case config.ErrorThreshold > 0 && duration >= config.ErrorThreshold:
		return ERROR
	case config.WarningThreshold > 0 && duration >= config.WarningThreshold:
		return WARN
	default:
		return INFO
	}
}

// getLogLevelColor returns ANSI color codes for different log levels
func getLogLevelColor(level LogLevel) string {
	switch level {
	case ERROR:
		return "\033[31m" // Red
	case WARN:
		return "\033[33m" // Yellow
	default:
		return "\033[32m" // Green
	}
}

// New returns a gin middleware for logging response times
func New(config Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip if middleware is disabled
		if !config.Enabled {
			c.Next()
			return
		}

		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Determine log level
		logLevel := determineLogLevel(duration, config)

		// Skip logging if it's not a slow request and we're not logging all requests
		if !config.LogAllRequests && logLevel == INFO {
			return
		}

		// Get status code
		statusCode := c.Writer.Status()

		// Create log message with color coding
		reset := "\033[0m"
		colorCode := getLogLevelColor(logLevel)

		message := fmt.Sprintf("%s[%s]%s %v | %3d | %13v | %15s | %-7s %s",
			colorCode,
			logLevel,
			reset,
			time.Now().Format("2006/01/02 - 15:04:05"),
			statusCode,
			duration,
			c.ClientIP(),
			c.Request.Method,
			c.Request.URL.Path,
		)

		// Add error if present
		if len(c.Errors) > 0 {
			message += " | " + c.Errors.String()
		}

		// Log using configured logger
		if config.Logger != nil {
			config.Logger(message)
		} else {
			log.Print(message)
		}
	}
}

// Enable enables the middleware
func (c *Config) Enable() {
	c.Enabled = true
}

// Disable disables the middleware
func (c *Config) Disable() {
	c.Enabled = false
}
