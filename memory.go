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
	frequency := 60 // 60 seconds
	if len(ctx.Query("frequency")) > 0 {
		frequency, _ = strconv.Atoi(ctx.Query("frequency"))
	}
	maxMemory := 1024 * 1 // 1 GB
	if len(ctx.Query("max-memory")) > 0 {
		maxMemory, _ = strconv.Atoi(ctx.Query("max-memory"))
	}
	slice := make([]byte, 1024*1024)
	go func() {
		for {
			slice = append(slice, slice...)
			memStats := runtime.MemStats{}
			runtime.ReadMemStats(&memStats)
			fmt.Printf("Memory usage: %d MB\n", memStats.Alloc/1024/1024)
			if memStats.Alloc/1024/1024 > uint64(maxMemory) {
				break
			}
			time.Sleep(time.Second * time.Duration(frequency))
		}
	}()
	output := "Memory leak simulation started"
	log.Println(output)
	ctx.String(http.StatusOK, output)
}
