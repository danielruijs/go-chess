package main

import (
	"go-chess/internal/server"
	"log"
	"net/http"
)

func main() {
	config := parseFlags()

	metrics := server.NewMetrics()
	matchmaker := server.NewMatchmaker(metrics)
	go matchmaker.Run()

	webSocketHandler := server.NewWebSocketHandler(matchmaker, config.AllowLocalhost)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", metrics.MetricsHandler())

	http.HandleFunc("/ws", webSocketHandler.ServeWS)

	log.Print("Started server")
	go func() {
		log.Print("Started metrics server on localhost:2115")
		log.Fatal(http.ListenAndServe("localhost:2115", metricsMux))
	}()
	log.Fatal(http.ListenAndServe("localhost:8085", nil))
}
