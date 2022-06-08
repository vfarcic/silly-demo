package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/", rootHandler)
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	r.Run(fmt.Sprintf(":%s", port))
}

func rootHandler(c *gin.Context) {
	code := http.StatusOK
	output := "This is a silly demo"
	if len(c.Query("fail")) > 0 {
		code = http.StatusInternalServerError
		output = "Something terrible happened"
	}
	version := os.Getenv("VERSION")
	if len(version) > 0 {
		output = fmt.Sprintf("%s version %s", output, version)
	}
	if len(c.Query("html")) > 0 {
		output = fmt.Sprintf("<h1>%s</h1>", output)
	}
	output = fmt.Sprintf("%s\n", output)
	c.String(code, output)
}
