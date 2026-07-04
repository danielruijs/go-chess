package main

import (
	"flag"
	"fmt"
	"strings"
)

type Config struct {
	AllowedOrigins []string
	CookieDomain   string
	DatabaseDSN    string

	Port        int
	MetricsPort int
}

func parseFlags() (Config, error) {
	var cfg Config
	var origins string

	flag.StringVar(&origins, "allowed-origins", "", "comma-separated list of allowed CORS/WebSocket origins")
	flag.StringVar(&cfg.CookieDomain, "cookie-domain", "", "domain for session cookies")
	flag.StringVar(&cfg.DatabaseDSN, "db-dsn", "", "PostgreSQL connection DSN")
	flag.IntVar(&cfg.Port, "port", 8085, "port for gameplay traffic")
	flag.IntVar(&cfg.MetricsPort, "metrics-port", 2115, "port for metrics server")

	flag.Parse()

	if origins == "" {
		return Config{}, fmt.Errorf("-allowed-origins is required")
	}
	if cfg.CookieDomain == "" {
		return Config{}, fmt.Errorf("-cookie-domain is required")
	}
	if cfg.DatabaseDSN == "" {
		return Config{}, fmt.Errorf("-db-dsn is required")
	}

	if cfg.Port < 1 || cfg.Port > 65535 {
		return Config{}, fmt.Errorf("invalid -port: %d (must be 1-65535)", cfg.Port)
	}
	if cfg.MetricsPort < 1 || cfg.MetricsPort > 65535 {
		return Config{}, fmt.Errorf("invalid -metrics-port: %d (must be 1-65535)", cfg.MetricsPort)
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
