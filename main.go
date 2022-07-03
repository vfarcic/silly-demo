package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/fibonacci", fibonacciHandler)
	r.GET("/", rootHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func rootHandler(c *gin.Context) {
	if len(c.Query("fail")) > 0 {
		c.String(http.StatusBadRequest, "Something terrible happened")
		return
	}
	version := os.Getenv("VERSION")
	output := "This is a silly demo"
	if len(version) > 0 {
		output = fmt.Sprintf("%s version %s", output, version)
	}
	if len(c.Query("html")) > 0 {
		output = fmt.Sprintf("<h1>%s</h1>", output)
	}
	output = fmt.Sprintf("%s\n", output)
	c.String(http.StatusOK, output)
}

// Handle the /fibonacci path
func fibonacciHandler(c *gin.Context) {
	// Get number parameter from query
	number, err := strconv.Atoi(c.Query("number"))
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	fib := calculateFibonacci(number)
	c.String(http.StatusOK, fmt.Sprintf("%d", fib))
}

func calculateFibonacci(n int) int {
	if n <= 1 {
		return n
	}
	return calculateFibonacci(n-1) + calculateFibonacci(n-2)
}
