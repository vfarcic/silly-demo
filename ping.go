package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

func pingHandler(ctx *gin.Context) {
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	req := resty.New().R().SetHeaderMultiValues(ctx.Request.Header).SetHeader("Content-Type", "application/text")
	url := ctx.Query("url")
	if len(url) == 0 {
		url = os.Getenv("PING_URL")
		if len(url) == 0 {
			httpErrorBadRequest(errors.New("url is empty"), ctx)
			return
		}
	}
	slog.Info("Sending a ping", "URL", url)
	resp, err := req.Get(url)
	if err != nil {
		httpErrorBadRequest(err, ctx)
		return
	}
	slog.Info(resp.String())
	ctx.String(http.StatusOK, resp.String())
}
