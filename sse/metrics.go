package sse

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"runtime"
	"time"
)

type SystemMetrics struct {
	Time       string  `json:"time"`
	Goroutines int     `json:"goroutines"`
	AllocMB    float64 `json:"alloc_mb"`
}


func MetricsStreamHandler(w http.ResponseWriter, r *http.Request) {

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}


	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() 
	for {
		select {

		case <-r.Context().Done():
			fmt.Println("SSE: Клієнт відключився, потік зупинено.")
			return

		case <-ticker.C:
			var memStats runtime.MemStats
			runtime.ReadMemStats(&memStats)

			allocMB := math.Round((float64(memStats.Alloc)/1024/1024)*100) / 100

			metrics := SystemMetrics{
				Time:       time.Now().Format(time.RFC3339),
				Goroutines: runtime.NumGoroutine(),
				AllocMB:    allocMB,
			}

			jsonData, err := json.Marshal(metrics)
			if err != nil {
				continue
			}


			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush() 
		}
	}
}