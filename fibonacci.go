package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RespondFibonacci(message string) string {
	number, err := strconv.Atoi(message)
	if err != nil {
		return fmt.Sprintf("%s is not a number", message)
	}
	return strconv.Itoa(calculateFibonacci(number))
}

func fibonacciHandler(ctx *gin.Context) {
	slog.Debug("Handling request", "URI", ctx.Request.RequestURI)
	number, err := strconv.Atoi(ctx.Query("number"))
	if err != nil {
		ctx.String(http.StatusBadRequest, err.Error())
		return
	}
	fib := calculateFibonacci(number)
	if os.Getenv("PUBLISH") == "nats" {
		NatsPublish(fmt.Sprintf("Fibonacci %d = %d", number, fib))
	}
	ctx.String(http.StatusOK, fmt.Sprintf("%d", fib))
}

func calculateFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return calculateFibonacci(n-1) + calculateFibonacci(n-2)
}
