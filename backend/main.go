package main

import (
	"go-chess/internal/auth"
	"go-chess/internal/server"
	"log"
	"net/http"
)

func main() {
	config, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	metrics := server.NewMetrics()
	matchmaker := server.NewMatchmaker(metrics)
	go matchmaker.Run()

	userStore := auth.NewUserStore()
	sessionStore := auth.NewSessionStore()
	authHandler := auth.NewAuthHandler(userStore, sessionStore, config.AllowedOrigins, config.CookieDomain)

	webSocketHandler := server.NewWebSocketHandler(matchmaker, config.AllowedOrigins, sessionStore)
	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", metrics.MetricsHandler())

	http.HandleFunc("/ws", webSocketHandler.ServeWS)

	http.HandleFunc("/api/register", authHandler.CORSMiddleware(authHandler.Register))
	http.HandleFunc("/api/login", authHandler.CORSMiddleware(authHandler.Login))
	http.HandleFunc("/api/logout", authHandler.CORSMiddleware(authHandler.Logout))
	http.HandleFunc("/api/check", authHandler.CORSMiddleware(authHandler.CheckAuth))

	log.Print("Started server")
	go func() {
		log.Print("Started metrics server on :2115")
		log.Fatal(http.ListenAndServe(":2115", metricsMux))
	}()
	log.Fatal(http.ListenAndServe(":8085", nil))
}
