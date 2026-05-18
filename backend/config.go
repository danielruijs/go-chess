package main

import "flag"

type Config struct {
	AllowLocalhost bool
}

func parseFlags() Config {
	dev := flag.Bool("dev", false, "enable dev features (allow localhost origins)")

	flag.Parse()
	return Config{
		AllowLocalhost: *dev,
	}
}
