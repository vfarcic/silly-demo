package main

import (
	"fmt"
	"log/slog"
	"net/http"
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
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fib := calculateFibonacci(number)
	ctx.JSON(http.StatusOK, gin.H{"number": number, "fibonacci": fib})
}

func calculateFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return calculateFibonacci(n-1) + calculateFibonacci(n-2)
}
