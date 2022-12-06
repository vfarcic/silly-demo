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
const version = "v0.0.1"

func main() {
	// OpenTelemetry
	l := log.New(os.Stdout, "", 0)
	f, err := os.Create("traces.txt")
	if err != nil {
		l.Fatal(err)
	}
	defer f.Close()
	tp, err := initTrace(f)
	if err != nil {
		l.Fatal(err)
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
