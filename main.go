package main

import (
	"log"
	"net/http"
	"pr-6-go/sse"
	"pr-6-go/ws"
)

func main() {
	hub := ws.NewHub()
	go hub.Run()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /api/v1/metrics/stream", sse.MetricsStreamHandler)

	mux.HandleFunc("GET /api/v1/ws", func(w http.ResponseWriter, r *http.Request) {
		ws.ServeWs(hub, w, r)
	})

	log.Println("Сервер запущено на http://localhost:8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Помилка запуску сервера:", err)
	}
}