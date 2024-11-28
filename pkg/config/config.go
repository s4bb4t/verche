package config

import (
	"os"
	"strings"
)

type Config struct {
	GoVersion string
}

func MustLoad() *Config {
	var cfg Config

	cfg.GoVersion = strings.TrimSpace(os.Getenv("GO_PROJECT_VERSION"))
	if cfg.GoVersion == "" {
		panic("GO_PROJECT_VERSION is not set")
	}

	return &cfg
}
