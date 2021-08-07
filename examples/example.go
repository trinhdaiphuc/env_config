package main

import (
	"fmt"
	_ "github.com/joho/godotenv/autoload"
	"github.com/trinhdaiphuc/env_config"
	"time"
)

type Config struct {
	Host    string  `env:"HOST,localhost"`
	Port    int     `env:"PORT,8080"`
	Bytes   []byte  `env:"ENV_BYTES,foo"`
	Float32 float32 `env:"ENV_FLOAT,12.34"`
	None    string
	Timeout time.Duration `env:"TIMEOUT"`
	Date    time.Time     `env:"ENV_DATE"`
}

func main() {
	cfg := &Config{}
	if err := env_config.EnvStruct(cfg); err != nil {
		panic(err)
	}
	fmt.Printf("Config %+v\n", cfg)
}
