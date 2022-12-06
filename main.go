package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
)

const name = "silly-demo"

func main() {
	// OpenTelemetry
	var err error
	tp, err = initTrace()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()
	otel.SetTracerProvider(tp)

	// Server
	r := gin.Default()
	r.Use(otelgin.Middleware(name))
	r.GET("/fibonacci", fibonacciHandler)
	r.POST("/slack", slackHandler)
	r.POST("/video", videoPostHandler)
	r.GET("/videos", videosGetHandler)
	r.GET("/", rootHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}
