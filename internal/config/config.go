package config

import "os"

type Config struct {
	HTTPPort string
}

func LoadConfig() Config {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	return Config{HTTPPort: httpPort}
}
