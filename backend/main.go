package main

import (
	"context"
	"fmt"
	"go-chess/internal/auth"
	"go-chess/internal/db"
	"go-chess/internal/db/sqlc"
	"go-chess/internal/server"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	config, err := parseFlags()
	if err != nil {
		log.Fatal(err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	metrics := server.NewMetrics()
	matchmaker, err := server.NewMatchmaker(metrics)
	if err != nil {
		log.Fatalf("failed to create matchmaker: %v", err)
	}
	go matchmaker.Run(ctx)

	pool, err := db.Connect(ctx, config.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	queries := sqlc.New(pool)
	userStore, err := auth.NewUserStore(queries)
	if err != nil {
		log.Fatalf("failed to create user store: %v", err)
	}
	sessionStore, err := auth.NewSessionStore(ctx)
	if err != nil {
		log.Fatalf("failed to create session store: %v", err)
	}
	authHandler := auth.NewAuthHandler(userStore, sessionStore, config.AllowedOrigins, config.CookieDomain)

	webSocketHandler, err := server.NewWebSocketHandler(ctx, matchmaker, config.AllowedOrigins, sessionStore)
	if err != nil {
		log.Fatalf("failed to create websocket handler: %v", err)
	}

	metricsMux := http.NewServeMux()
	metricsMux.Handle("/metrics", metrics.MetricsHandler())

	http.HandleFunc("/ws", webSocketHandler.ServeWS)

	http.HandleFunc("/api/register", authHandler.CORSMiddleware(authHandler.Register))
	http.HandleFunc("/api/login", authHandler.CORSMiddleware(authHandler.Login))
	http.HandleFunc("/api/logout", authHandler.CORSMiddleware(authHandler.Logout))
	http.HandleFunc("/api/check", authHandler.CORSMiddleware(authHandler.CheckAuth))

	srv := &http.Server{Addr: fmt.Sprintf(":%d", config.Port)}
	metricsSrv := &http.Server{Addr: fmt.Sprintf(":%d", config.MetricsPort), Handler: metricsMux}

	go func() {
		log.Printf("started metrics server on :%d", config.MetricsPort)
		if err := metricsSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("metrics server failed: %v", err)
		}
	}()

	go func() {
		log.Printf("started server on :%d", config.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()

	<-ctx.Done()
	log.Println("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("server shutdown error: %v", err)
	}
	if err := metricsSrv.Shutdown(shutdownCtx); err != nil {
		log.Printf("metrics server shutdown error: %v", err)
	}

	log.Println("server exited")
}
