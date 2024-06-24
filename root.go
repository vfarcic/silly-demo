package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/open-feature/go-sdk/openfeature"
)

func rootHandler(ctx *gin.Context) {
	// OpenFeature
	messageFlag, _ := OpenFeatureClient.BooleanValue(
		context.Background(), "root-message", false, openfeature.EvaluationContext{},
	)

	if len(ctx.Query("fail")) > 0 || len(os.Getenv("FAIL")) > 0 {
		ctx.String(http.StatusInternalServerError, "Something terrible happened")
		return
	}
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	version := os.Getenv("VERSION")
	output := os.Getenv("MESSAGE")
	if len(output) == 0 {
		if messageFlag {
			output = "This is a silly demo with OpenFeature"
		} else {
			output = "This is a silly demo"
		}
	}
	if len(version) > 0 {
		output = fmt.Sprintf("%s version %s", output, version)
	}
	if len(ctx.Query("html")) > 0 {
		output = fmt.Sprintf("<h1>%s</h1>", output)
	}
	output = fmt.Sprintf("%s\n", output)
	ctx.String(http.StatusOK, output)
}
