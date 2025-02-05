package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"silly-demo/internal/service"

	"github.com/gin-gonic/gin"
)

func FibonacciHandler(ctx *gin.Context) {
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	number, err := strconv.Atoi(ctx.Query("number"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fib := service.CalculateFibonacci(number)
	ctx.JSON(http.StatusOK, gin.H{"number": number, "fibonacci": fib})
}
