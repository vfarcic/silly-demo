package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var memoryLeakSlice = make([]byte, 1024*1024)

func memoryLeakHandler(ctx *gin.Context) {
	maxMemory := 0
	if len(ctx.Query("max-memory")) > 0 {
		maxMemory, _ = strconv.Atoi(ctx.Query("max-memory"))
	}
	frequency := 60 // 60 seconds
	if len(ctx.Query("frequency")) > 0 {
		frequency, _ = strconv.Atoi(ctx.Query("frequency"))
	}
	go func() {
		memoryLeak(maxMemory, frequency)
	}()
	ctx.String(http.StatusOK, "Memory leak simulation started")
}

func memoryLeak(maxMemory, frequency int) {
	if maxMemory <= 0 {
		maxMemory = 1024 * 1 // 1 GB
		if len(os.Getenv("MEMORY_LEAK_MAX_MEMORY")) > 0 {
			maxMemory, _ = strconv.Atoi(os.Getenv("MEMORY_LEAK_MAX_MEMORY"))
		}
	}
	if frequency <= 0 {
		frequency = 60
		if len(os.Getenv("MEMORY_LEAK_FREQUENCY")) > 0 {
			frequency, _ = strconv.Atoi(os.Getenv("MEMORY_LEAK_FREQUENCY"))
		}
	}
	log.Printf(
		"Memory leak simulation started with max memory set to %d MB and frequency set to %d seconds",
		maxMemory,
		frequency,
	)
	for {
		memoryLeakSlice = append(memoryLeakSlice, memoryLeakSlice...)
		memStats := runtime.MemStats{}
		runtime.ReadMemStats(&memStats)
		fmt.Printf("Memory usage: %d MB\n", memStats.Alloc/1024/1024)
		if maxMemory > 0 && memStats.Alloc/1024/1024 > uint64(maxMemory) {
			log.Println("Memory leak simulation ended")
			break
		}
		time.Sleep(time.Second * time.Duration(frequency))
	}
}
