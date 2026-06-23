package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	AllowedOrigins []string
	CookieDomain   string

	Port        string
	MetricsPort string
}

func parseFlags() (Config, error) {
	var cfg Config
	var origins string

	flag.StringVar(&origins, "allowed-origins", "", "comma-separated list of allowed CORS/WebSocket origins")
	flag.StringVar(&cfg.CookieDomain, "cookie-domain", "", "domain for session cookies")
	flag.StringVar(&cfg.Port, "port", "8085", "port for gameplay traffic")
	flag.StringVar(&cfg.MetricsPort, "metrics-port", "2115", "port for metrics server")

	flag.Parse()

	if origins == "" {
		return Config{}, fmt.Errorf("-allowed-origins is required")
	}
	if cfg.CookieDomain == "" {
		return Config{}, fmt.Errorf("-cookie-domain is required")
	}

	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}
	if !strings.HasPrefix(cfg.MetricsPort, ":") {
		cfg.MetricsPort = ":" + cfg.MetricsPort
	}

	parts := strings.SplitSeq(origins, ",")
	for p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			cfg.AllowedOrigins = append(cfg.AllowedOrigins, trimmed)
		}
	}

	if len(cfg.AllowedOrigins) == 0 {
		return Config{}, fmt.Errorf("no valid non-empty origins parsed from -allowed-origins")
	}

	return cfg, nil
}
