package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"github.com/open-feature/go-sdk/openfeature"
)

var OpenFeatureClient = openfeature.NewClient("silly-demo")

func init() {
	err := openfeature.SetProviderAndWait(flagd.NewProvider())
	if err != nil {
		log.Printf("Failed to set the OpenFeature provider: %v", err)
	}
}

func main() {
	// Logging
	log.SetOutput(os.Stderr)
	if os.Getenv("MEMORY_LEAK_MAX_MEMORY") != "" {
		go func() { memoryLeak(0, 0) }()
	}

	// Server
	log.Println("Starting server...")
	router := gin.New()
	router.GET("/fibonacci", fibonacciHandler)
	router.POST("/video", videoPostHandler)
	router.GET("/videos", videosGetHandler)
	router.GET("/ping", pingHandler)
	router.GET("/memory-leak", memoryLeakHandler)
	router.GET("/", rootHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	router.Run(fmt.Sprintf(":%s", port))
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
