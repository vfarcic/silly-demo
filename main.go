package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"silly-demo/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var OpenFeatureClient *openfeature.Client

func init() {
	provider, err := flagd.NewProvider()
	if err != nil {
		log.Fatalf("failed to create flagd provider: %v", err)
	}
	openfeature.SetProviderAndWait(provider)
	OpenFeatureClient = openfeature.NewClient("silly-demo")
}

func main() {
	log.SetOutput(os.Stderr)
	if os.Getenv("DEBUG") == "true" {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	if os.Getenv("MEMORY_LEAK_MAX_MEMORY") != "" {
		go func() { memoryLeak(0, 0) }()
	}

	NatsSubscribe()
	go NatsPublishLoop()

	// Server
	log.Println("Starting server...")
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.GET("/fibonacci", handlers.FibonacciHandler)
	router.POST("/video", handlers.VideoPostHandler)
	router.GET("/videos", handlers.VideosGetHandler)
	router.GET("/ping", handlers.PingHandler)
	router.GET("/memory-leak", memoryLeakHandler)
	router.GET("/", rootHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", port),
		Handler: router.Handler(),
	}

	// Signals
	if len(os.Getenv("NO_SIGNALS")) > 0 {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("HTTP server error: %v", err)
		}
	} else {
		go func() {
			if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
				log.Fatalf("HTTP server error: %v", err)
			}
			log.Println("Stopped serving new connections.")
		}()
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 60*time.Second)
		defer shutdownRelease()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Fatalf("HTTP shutdown error: %v", err)
		}
		openfeature.Shutdown()
		log.Println("Graceful shutdown complete.")
	}
	// <-ctx.Done()
}

func httpErrorBadRequest(err error, ctx *gin.Context) {
	httpError(err, ctx, http.StatusBadRequest)
}

func httpErrorInternalServerError(err error, ctx *gin.Context) {
	httpError(err, ctx, http.StatusInternalServerError)
}

func httpError(err error, ctx *gin.Context, status int) {
	log.Println(err.Error())
	ctx.String(status, err.Error())
}
