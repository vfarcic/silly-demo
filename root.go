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
	if len(ctx.Query("fail")) > 0 || len(os.Getenv("FAIL")) > 0 {
		ctx.String(http.StatusInternalServerError, "Something terrible happened")
		return
	}
	// FIXME: Enable debugging with a flag
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	version := os.Getenv("VERSION")
	output, _ := OpenFeatureClient.StringValue(
		context.Background(), "output", "This is a silly demo", openfeature.EvaluationContext{},
	)
	includeVersion, _ := OpenFeatureClient.BooleanValue(
		context.Background(), "include-version", false, openfeature.EvaluationContext{},
	)
	if includeVersion && len(version) > 0 {
		output = fmt.Sprintf("%s version %s", output, version)
	}
	if len(ctx.Query("html")) > 0 {
		output = fmt.Sprintf("<h1>%s</h1>", output)
	}
	output = fmt.Sprintf("%s\n", output)
	ctx.String(http.StatusOK, output)
}
