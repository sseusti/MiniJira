package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	HTTPPort  string
	LogLevel  string
	LogFormat string
}

func LoadConfig() (Config, error) {
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	LogLevel := os.Getenv("LOG_LEVEL")
	if LogLevel == "" {
		LogLevel = "info"
	}
	err := Validate(LogLevel, allowedLevels)
	if err != nil {
		return Config{}, err
	}

	LogFormat := os.Getenv("LOG_FORMAT")
	if LogFormat == "" {
		LogFormat = "text"
	}
	err = Validate(LogLevel, allowedLevels)
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPPort:  httpPort,
		LogLevel:  LogLevel,
		LogFormat: LogFormat}, nil
}

var allowedLevels = map[string]struct{}{
	"debug": {},
	"info":  {},
	"warn":  {},
	"error": {},
}

var allowedFormats = map[string]struct{}{
	"text": {},
	"json": {},
}

func Validate(check string, allow map[string]struct{}) error {
	check = strings.TrimSpace(strings.ToLower(check))
	if _, ok := allow[check]; !ok {
		return fmt.Errorf("%s is not allowed", check)
	}
	return nil
}
