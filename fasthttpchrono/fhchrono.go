package fasthttpchrono

import (
	"fmt"
	"log"
	"time"

	"github.com/valyala/fasthttp"
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
	// Skip paths from logging (e.g., health checks)
	SkipPaths []string
}

// DefaultConfig returns a default configuration
func DefaultConfig() Config {
	return Config{
		Enabled:          true,
		WarningThreshold: 500 * time.Millisecond,  // 500ms default for warnings
		ErrorThreshold:   2000 * time.Millisecond, // 2s default for errors
		LogAllRequests:   false,                   // only log slow requests by default
		Logger:           log.Printf,              // use standard logger by default
		SkipPaths:        []string{},              // no paths skipped by default
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

// shouldSkipPath checks if the current path should be skipped from logging
func shouldSkipPath(path string, skipPaths []string) bool {
	for _, skipPath := range skipPaths {
		if skipPath == path {
			return true
		}
	}
	return false
}

// New returns a fasthttp middleware for logging response times
func New(config Config) func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(next fasthttp.RequestHandler) fasthttp.RequestHandler {
		return func(ctx *fasthttp.RequestCtx) {
			// Skip if middleware is disabled or path should be skipped
			if !config.Enabled || shouldSkipPath(string(ctx.Path()), config.SkipPaths) {
				next(ctx)
				return
			}

			// Start timer
			start := time.Now()

			// Process request
			next(ctx)

			// Calculate duration
			duration := time.Since(start)

			// Determine log level
			logLevel := determineLogLevel(duration, config)

			// Skip logging if it's not a slow request and we're not logging all requests
			if !config.LogAllRequests && logLevel == INFO {
				return
			}

			// Get status code
			statusCode := ctx.Response.StatusCode()

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
				ctx.RemoteIP().String(),
				string(ctx.Method()),
				string(ctx.Path()),
			)

			// Add user agent if available
			userAgent := string(ctx.UserAgent())
			if userAgent != "" {
				message += fmt.Sprintf(" | %s", userAgent)
			}

			// Log using configured logger
			if config.Logger != nil {
				config.Logger(message)
			} else {
				log.Print(message)
			}
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

// AddSkipPath adds a path to skip from logging
func (c *Config) AddSkipPath(path string) {
	c.SkipPaths = append(c.SkipPaths, path)
}
