package main

import (
	"log"
	"time"

	"github.com/ashraful88/chrono/fasthttpchrono"
	"github.com/valyala/fasthttp"
)

func main() {
	// Create config
	config := fasthttpchrono.DefaultConfig()

	// Customize settings
	config.WarningThreshold = 200 * time.Millisecond
	config.ErrorThreshold = 1 * time.Second
	config.LogAllRequests = true

	// Skip health check endpoint from logging
	config.AddSkipPath("/health")

	// Custom logger
	config.Logger = func(format string, v ...interface{}) {
		log.Printf("[FastHTTP] "+format, v...)
	}

	// Create handler
	handler := fasthttpchrono.New(config)(fastHTTPHandler)

	// Start server
	log.Fatal(fasthttp.ListenAndServe(":8080", handler))
}

func fastHTTPHandler(ctx *fasthttp.RequestCtx) {
	switch string(ctx.Path()) {
	case "/fast":
		ctx.SetStatusCode(200)
		ctx.SetBodyString("fast response")

	case "/slow":
		time.Sleep(600 * time.Millisecond) // Will trigger WARNING
		ctx.SetStatusCode(200)
		ctx.SetBodyString("slow response")

	case "/very-slow":
		time.Sleep(2100 * time.Millisecond) // Will trigger ERROR
		ctx.SetStatusCode(200)
		ctx.SetBodyString("very slow response")

	case "/health":
		ctx.SetStatusCode(200)
		ctx.SetBodyString("OK")

	default:
		ctx.SetStatusCode(404)
	}
}
