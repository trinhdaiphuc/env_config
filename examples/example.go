package main

import (
	"fmt"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/trinhdaiphuc/env_config"
)

type Config struct {
	Host      string        `env:"HOST;default=localhost"`
	Port      int           `env:"PORT;default=8080"`
	Bytes     []byte        `env:"ENV_BYTES;default=foo,bar"`
	Float32   float32       `env:"ENV_FLOAT;default=12.34"`
	Timeout   time.Duration `env:"TIMEOUT"`
	Date      time.Time     `env:"ENV_DATE"`
	Addresses []string      `env:"ADDRESSES"`
	None      string
}

func main() {
	cfg := &Config{}
	if err := env_config.LoadConfig(cfg); err != nil {
		panic(err)
	}
	fmt.Printf("Config %+v\n", cfg)
}
