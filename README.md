# env_config

`env_config` is a Go package designed to simplify the process of loading configuration values from environment variables into your application's structs. This package supports various data types, including custom and complex types, and allows for the use of tag options to customize the behavior of value parsing.

## Features

- Load configuration from environment variables directly into structs.
- Supports basic types like `string`, `int`, `uint`, `float`, `bool`, and their slices.
- Supports complex types like `time.Duration` and `time.Time`.
- Customizable tag options to control how environment variables are parsed and set.

## Installation

```bash
go get github.com/trinhdaiphuc/env_config
```

## Usage

### Basic Usage

Define your configuration struct and use the `env` tag to specify the environment variable names:

```go
package main

import (
	"fmt"
	"log"
	"github.com/trinhdaiphuc/env_config"
)

type Config struct {
	Port        int           `env:"PORT"`
	DatabaseURL string        `env:"DATABASE_URL"`
	Debug       bool          `env:"DEBUG"`
	Timeout     time.Duration `env:"TIMEOUT"`
}

func main() {
	var config Config
	if err := env_config.Load(&config); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%+v\n", config)
}
```

### Tag Options

You can customize the parsing behavior using tag options. For example, you can specify a default value or a delimiter for slice types.

```go
type Config struct {
	Port        int           `env:"PORT;default=8080"`
	DatabaseURL string        `env:"DATABASE_URL"`
	Debug       bool          `env:"DEBUG;default=false"`
	Timeout     time.Duration `env:"TIMEOUT;default=5s"`
	Hosts       []string      `env:"HOSTS;delimiter=,"`
}
```

### Supported Types

The package supports the following types:

- `string`
- `int`, `int8`, `int16`, `int32`, `int64`
- `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- `float32`, `float64`
- `bool`
- `[]string`
- `[]int`, `[]int8`, `[]int16`, `[]int32`, `[]int64`
- `[]uint`, `[]uint8`, `[]uint16`, `[]uint32`, `[]uint64`
- `[]float32`, `[]float64`
- `[]bool`
- `time.Duration`
- `time.Time` (parsed using `time.RFC3339` format)

## Custom Strategies

If you need to support custom types, you can implement the `TypeStrategy` interface and register your custom strategy.

```go
type CustomType struct {
	// Your fields here
}

type CustomTypeStrategy struct{}

func (s CustomTypeStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	// Your custom parsing logic here
}

func main() {
	env_config.RegisterStrategy(reflect.TypeOf(CustomType{}), CustomTypeStrategy{})
	// Now you can use CustomType in your config struct
}
```

## Error Handling

The `Load` function returns an error if any required environment variables are missing or if any values cannot be parsed. You can handle these errors as needed in your application.

## Testing

To test your configuration loading logic, you can set environment variables in your test cases and use the `Load` function as usual.

```go
func TestConfigLoading(t *testing.T) {
	os.Setenv("PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/db")
	os.Setenv("DEBUG", "true")
	os.Setenv("TIMEOUT", "10s")
	os.Setenv("HOSTS", "host1,host2,host3")

	var config Config
	if err := env_config.Load(&config); err != nil {
		t.Fatal(err)
	}

	// Add your assertions here
}
```

## Contributing

Contributions are welcome! Please open an issue or submit a pull request on GitHub.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
