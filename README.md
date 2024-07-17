# env_config

> Load environment config and utilities for get default environment variables which have or not.

## Install

```shell
go get -u github.com/trinhdaiphuc/env_config
```

## Usage

Define your config with struct tag `env:{env_variable,default_value}`. Ex:

```go
type Config struct {
    Host    string  `env:"HOST,localhost"`
    Port    int     `env:"PORT,8080"`
    Bytes   []byte  `env:"ENV_BYTES,foo"`
    Float32 float32 `env:"ENV_FLOAT"`
    None    string
    Timeout time.Duration `env:"TIMEOUT"`
    Date    time.Time     `env:"ENV_DATE"`
}
```

The library will get environment value from `env_variable` if value is empty it will add value with `default_value` you
describe. Example:

I use `.env` file for define environments easily:

```dotenv
HOST=127.0.0.1
PORT=8081
ENV_BYTES=foo
ENV_FLOAT=
TIMEOUT=5s
ENV_DATE=2021-08-07T15:04:05Z
```

I will get these environments into my config struct. [example.go](./examples/example.go)

```go
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
```

Now you can run `example.go` to check this

```shell
go run examples/example.go

Config &{127.0.0.1 8081 [102 111 111] 12.34  5s 2021-08-07 15:04:05 +0000 UTC}
```