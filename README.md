# env_config

> Load environment config and utilities for getting default environment variables, whether they exist or not.

## Install

```shell
go get -u github.com/trinhdaiphuc/env_config
```

## Usage

Define your config with struct tag `env:"{env_variable},default={default_value}"`. For example:

```go
type Config struct {
    Host        string        `env:"HOST,default=localhost"`
    Port        int           `env:"PORT,default=8080"`
    Bytes       []byte        `env:"ENV_BYTES,default=foo"`
    Float32     float32       `env:"ENV_FLOAT,default=12.34"`
    None        string
    Timeout     time.Duration `env:"TIMEOUT"`
    Date        time.Time     `env:"ENV_DATE"`
    StringSlice []string      `env:"STRING_SLICE,default=foo,bar"`
    FloatSlice  []float64     `env:"FLOAT_SLICE,default=1.1,2.2,3.3"`
}
```

The library will get the environment value from `env_variable`. If the value is empty, it will use the `default_value` you describe. Example:

Use a `.env` file to define environments easily:

```dotenv
HOST=127.0.0.1
PORT=8081
ENV_BYTES=foo
ENV_FLOAT=
TIMEOUT=5s
ENV_DATE=2021-08-07T15:04:05Z
STRING_SLICE=foo,bar,baz
FLOAT_SLICE=1.1,2.2,3.3
```

You can load these environments into your config struct. [example.go](./examples/example.go):

```go
package main

import (
    "fmt"
    _ "github.com/joho/godotenv/autoload"
    "github.com/trinhdaiphuc/env_config"
    "time"
)

type Config struct {
    Host        string        `env:"HOST,default=localhost"`
    Port        int           `env:"PORT,default=8080"`
    Bytes       []byte        `env:"ENV_BYTES,default=foo"`
    Float32     float32       `env:"ENV_FLOAT,default=12.34"`
    None        string
    Timeout     time.Duration `env:"TIMEOUT"`
    Date        time.Time     `env:"ENV_DATE"`
    StringSlice []string      `env:"STRING_SLICE,default=foo,bar"`
    FloatSlice  []float64     `env:"FLOAT_SLICE,default=1.1,2.2,3.3"`
}

func main() {
    cfg := &Config{}
    if err := env_config.EnvStruct(cfg); err != nil {
        panic(err)
    }
    fmt.Printf("Config %+v\n", cfg)
}
```

Now you can run `example.go` to check this:

```shell
go run examples/example.go

Config &{127.0.0.1 8081 [102 111 111] 12.34  5s 2021-08-07 15:04:05 +0000 UTC [foo bar baz] [1.1 2.2 3.3]}
```
