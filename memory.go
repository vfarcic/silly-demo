package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func memoryLeakHandler(ctx *gin.Context) {
	frequencyString := ctx.Query("frequency")
	if len(frequencyString) == 0 {
		frequencyString = "60"
	}
	frequency, _ := strconv.Atoi(frequencyString)
	slice := make([]byte, 1024*1024)
	go func() {
		for {
			slice = append(slice, slice...)
			memStats := runtime.MemStats{}
			runtime.ReadMemStats(&memStats)
			fmt.Printf("Memory usage: %d MB\n", memStats.Alloc/1024/1024)
			time.Sleep(time.Second * time.Duration(frequency))
		}
	}()
	output := "Memory leak simulation started"
	log.Println(output)
	ctx.String(http.StatusOK, output)
}
