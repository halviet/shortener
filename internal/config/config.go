package config

import (
	"flag"
	"fmt"
	"github.com/caarlos0/env/v11"
)

type Config struct {
	Addr     string `env:"SERVER_ADDRESS"`
	BaseAddr string `env:"BASE_URL"`
}

func New() (Config, error) {
	cfg := parseFlags()
	envs, err := parseEnv()
	if err != nil {
		return Config{}, fmt.Errorf("config.New: %v", err)
	}

	// TODO: Could be done gracefully with reflect pkg
	// Env variables have the most priority,
	// if env not set, then it will use flag value or its default
	if envs.Addr != "" {
		cfg.Addr = envs.Addr
	}
	if envs.BaseAddr != "" {
		cfg.BaseAddr = envs.BaseAddr
	}

	// adding a trailing slash at the end of a base address
	if cfg.BaseAddr[len(cfg.BaseAddr)-1:] != "/" {
		cfg.BaseAddr += "/"
	}

	return cfg, nil
}

func parseFlags() Config {
	addr := flag.String("a", "localhost:8080", "address for HTTP-server (addr:port); default localhost:8080")
	baseAddr := flag.String("b", "", "sets base address for all resulting short urls; if not set uses -a flag address")

	flag.Parse()

	cfg := Config{
		Addr:     *addr,
		BaseAddr: *baseAddr,
	}

	return cfg
}

func parseEnv() (Config, error) {
	var cfg Config

	err := env.Parse(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf("parseEnv: %v", err)
	}

	// checking is -b flag was set
	// if it's not, then server address used as base
	if cfg.BaseAddr == "" {
		cfg.BaseAddr = "http://" + cfg.Addr + "/"
	}

	return cfg, nil
}
