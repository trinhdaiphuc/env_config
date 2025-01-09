package env_config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Config struct {
	BuildInConfig              *BuildInConfig              `env:"BUILD_IN"`
	ComplexConfig              *ComplexConfig              `env:"COMPLEX"`
	Scheduler                  *Scheduler                  `env:"SCHEDULER"`
	ServerConfig               *ServerConfig               `env:"SERVER"`
	NestedConfig               *NestedConfig               `env:"APP"`
	UninitializedPointerConfig *UninitializedPointerConfig `env:"UN_INIT"`
	EnvVarTaggedConfig         *EnvVarTaggedConfig         `env:"TAG"`
	Logger                     *Logger                     `env:"LOGGER"`
	NotSet                     *NotSetConfig               `env:"NOT_SET"`
	UnFieldConfig              string
}

type BuildInConfig struct {
	BoolValue    bool    `env:"BOOL_VALUE"`
	Float32Value float32 `env:"FLOAT32_VALUE"`
	Float64Value float64 `env:"FLOAT64_VALUE"`
	IntValue     int     `env:"INT_VALUE"`
	StringValue  string  `env:"STRING_VALUE"`
}

type ComplexConfig struct {
	BoolArray                      []bool        `env:"BOOL_ARRAY;delimiter=,"`
	Float32Array                   []float32     `env:"FLOAT32_ARRAY;delimiter= "`
	Float64Array                   []float64     `env:"FLOAT64_ARRAY;delimiter=-"`
	IntArray                       []int         `env:"INT_ARRAY;delimiter=_"`
	Int8Array                      []int8        `env:"INT8_ARRAY;delimiter=|"`
	Int16Array                     []int16       `env:"INT16_ARRAY;delimiter=|"`
	Int32Array                     []int32       `env:"INT32_ARRAY;delimiter=|"`
	Int64Array                     []int64       `env:"INT64_ARRAY;delimiter=|"`
	UintArray                      []uint        `env:"UINT_ARRAY;delimiter=|"`
	Uint8Array                     []uint8       `env:"UINT8_ARRAY;delimiter=|"`
	Uint16Array                    []uint16      `env:"UINT16_ARRAY;delimiter=|"`
	Uint32Array                    []uint32      `env:"UINT32_ARRAY;delimiter=|"`
	Uint64Array                    []uint64      `env:"UINT64_ARRAY;delimiter=|"`
	StringArray                    []string      `env:"STRING_ARRAY;delimiter=|"`
	StringArrayDefault             []string      `env:"STRING_ARRAY_DEFAULT;default=a,b,c,d"`
	StringArrayDefaultDelimiter    []string      `env:"STRING_ARRAY_DEFAULT_DELIMITER;default=a b c d;delimiter= "`
	StringArrayDefaultDelimiterEnv []string      `env:"STRING_ARRAY_DEFAULT_DELIMITER_ENV;default=a b c d;delimiter= "`
	Duration                       time.Duration `env:"DURATION"`
	Time                           time.Time     `env:"TIME"`
}

type EmptyConfig struct {
}

type ServerConfig struct {
	CacheConfig *RedisConfig `env:"REDIS"`
}

type NestedConfig struct {
	Database *DatabaseConfig `env:"DB"`
}

type UninitializedPointerConfig struct {
	LogLevel *string `env:"LOG_LEVEL"`
	Timeout  *int    `env:"TIMEOUT"`
}

type EnvVarTaggedConfig struct {
	AppName string  `env:"APP_NAME"`
	Debug   bool    `env:"DEBUG"`
	Pi      float64 `env:"PI"`
	Number  uint    `env:"NUMBER"`
}

type Logger struct {
	Level   string `env:"LEVEL"`
	Encoder string `env:"ENCODER"`
}

type NotSetConfig struct {
	Name          string   `env:"NAME"`
	Float32       float32  `env:"FLOAT32"`
	Float64       float64  `env:"FLOAT64"`
	Bool          bool     `env:"BOOL"`
	Int           int      `env:"INT"`
	Uint          uint     `env:"UINT"`
	String        string   `env:"STRING"`
	StringSlice   []string `env:"STRING_SLICE"`
	DefaultInt    int      `env:"DEFAULT_INT;default=10"`
	DefaultString string   `env:"DEFAULT_STRING;default=hello"`
	DefaultUint   uint     `env:"DEFAULT_UINT;default=10"`
	DefaultBool   bool     `env:"DEFAULT_BOOL;default=true"`
	DefaultFloat  float64  `env:"DEFAULT_FLOAT;default=3.14"`
}

type Scheduler struct {
	TimeInterval time.Duration `env:"INTERVAL"`
	StartAt      time.Time     `env:"START_AT"`
}

type RedisConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Password string `env:"PASSWORD"`
}

type DatabaseConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Username string `env:"USER"`
	Password string `env:"PASS"`
}

func TestLoadConfig(t *testing.T) {
	type args struct {
		cfg interface{}
	}

	var (
		timeNow = time.Now()
	)
	clean := setEnv(timeNow)
	defer clean()

	tests := []struct {
		name    string
		args    args
		wantCfg interface{}
		wantErr bool
	}{
		{
			name: "Test LoadConfig",
			args: args{
				cfg: &Config{},
			},
			wantCfg: &Config{
				BuildInConfig: &BuildInConfig{
					BoolValue:    true,
					Float32Value: 3.14,
					Float64Value: 3.14,
					IntValue:     10,
					StringValue:  "hello",
				},
				ComplexConfig: &ComplexConfig{
					BoolArray:                      []bool{true, false, false},
					Float32Array:                   []float32{3.14, 6.28, 0.3},
					Float64Array:                   []float64{3.14, 6.28, 12.2},
					IntArray:                       []int{1, 2, 3, 4, 5},
					Int8Array:                      []int8{1, 2, 3, 4, 5},
					Int16Array:                     []int16{1, 2, 3, 4, 5},
					Int32Array:                     []int32{1, 2, 3, 4, 5},
					Int64Array:                     []int64{1, 2, 3, 4, 5},
					UintArray:                      []uint{1, 2, 3, 4, 5},
					Uint8Array:                     []uint8{1, 2, 3, 4, 5},
					Uint16Array:                    []uint16{1, 2, 3, 4, 5},
					Uint32Array:                    []uint32{1, 2, 3, 4, 5},
					Uint64Array:                    []uint64{1, 2, 3, 4, 5},
					StringArray:                    []string{"a", "b", "c", "d"},
					StringArrayDefault:             []string{"a", "b", "c", "d"},
					StringArrayDefaultDelimiter:    []string{"a", "b", "c", "d"},
					StringArrayDefaultDelimiterEnv: []string{"x", "y", "z", "t"},
					Duration:                       5 * time.Minute,
					Time:                           timeNow.Truncate(time.Second),
				},
				Scheduler: &Scheduler{
					TimeInterval: 5 * time.Minute,
					StartAt:      timeNow.Truncate(time.Second),
				},
				ServerConfig: &ServerConfig{
					CacheConfig: &RedisConfig{
						Host:     "127.0.0.1",
						Port:     6379,
						Password: "secret",
					},
				},
				NestedConfig: &NestedConfig{
					Database: &DatabaseConfig{
						Host:     "localhost",
						Port:     3306,
						Username: "root",
						Password: "password",
					},
				},
				UninitializedPointerConfig: &UninitializedPointerConfig{
					LogLevel: stringPointer("info"),
					Timeout:  intPointer(10),
				},
				EnvVarTaggedConfig: &EnvVarTaggedConfig{
					AppName: "env_config",
					Debug:   true,
					Pi:      3.14,
					Number:  0,
				},
				Logger: &Logger{
					Level:   "debug",
					Encoder: "json",
				},
				UnFieldConfig: "",
				NotSet: &NotSetConfig{
					Name:          "",
					Float32:       0,
					Float64:       0,
					Bool:          false,
					Int:           0,
					Uint:          0,
					String:        "",
					StringSlice:   nil,
					DefaultInt:    10,
					DefaultString: "hello",
					DefaultUint:   10,
					DefaultBool:   true,
					DefaultFloat:  3.14,
				},
			},
			wantErr: false,
		},
		{
			name: "Test LoadConfig with empty config",
			args: args{
				cfg: &EmptyConfig{},
			},
			wantCfg: &EmptyConfig{

			},
			wantErr: false,
		},
		{
			name: "Test LoadConfig error",
			args: args{
				cfg: nil,
			},
			wantCfg: nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			assert.Equalf(t, tt.args.cfg, tt.wantCfg, "LoadConfig() = %v, want %v", toJSONStr(tt.args.cfg), toJSONStr(tt.wantCfg))
		})
	}
}

func setEnv(timeNow time.Time) func() {
	envMap := map[string]interface{}{
		"BUILD_IN_BOOL_VALUE":                        true,
		"BUILD_IN_FLOAT32_VALUE":                     3.14,
		"BUILD_IN_FLOAT64_VALUE":                     3.14,
		"BUILD_IN_INT_VALUE":                         10,
		"BUILD_IN_STRING_VALUE":                      "hello",
		"COMPLEX_BOOL_ARRAY":                         "true,false,false",
		"COMPLEX_FLOAT32_ARRAY":                      "3.14 6.28 0.3",
		"COMPLEX_FLOAT64_ARRAY":                      "3.14-6.28-12.2",
		"COMPLEX_INT_ARRAY":                          "1_2_3_4_5",
		"COMPLEX_INT8_ARRAY":                         "1|2|3|4|5",
		"COMPLEX_INT16_ARRAY":                        "1|2|3|4|5",
		"COMPLEX_INT32_ARRAY":                        "1|2|3|4|5",
		"COMPLEX_INT64_ARRAY":                        "1|2|3|4|5",
		"COMPLEX_UINT_ARRAY":                         "1|2|3|4|5",
		"COMPLEX_UINT8_ARRAY":                        "1|2|3|4|5",
		"COMPLEX_UINT16_ARRAY":                       "1|2|3|4|5",
		"COMPLEX_UINT32_ARRAY":                       "1|2|3|4|5",
		"COMPLEX_UINT64_ARRAY":                       "1|2|3|4|5",
		"COMPLEX_STRING_ARRAY":                       "a|b|c|d",
		"COMPLEX_DURATION":                           "5m",
		"COMPLEX_TIME":                               timeNow.Format(time.RFC3339),
		"COMPLEX_STRING_ARRAY_DEFAULT_DELIMITER_ENV": "x y z t",
		"SCHEDULER_INTERVAL":                         "5m",
		"SCHEDULER_START_AT":                         timeNow.Format(time.RFC3339),
		"SERVER_REDIS_HOST":                          "127.0.0.1",
		"SERVER_REDIS_PORT":                          6379,
		"SERVER_REDIS_PASSWORD":                      "secret",
		"APP_DB_HOST":                                "localhost",
		"APP_DB_PORT":                                3306,
		"APP_DB_USER":                                "root",
		"APP_DB_PASS":                                "password",
		"UN_INIT_LOG_LEVEL":                          "info",
		"UN_INIT_TIMEOUT":                            10,
		"TAG_APP_NAME":                               "env_config",
		"TAG_DEBUG":                                  true,
		"TAG_PI":                                     3.14,
		"LOGGER_LEVEL":                               "debug",
		"LOGGER_ENCODER":                             "json",
	}

	for env, val := range envMap {
		os.Setenv(env, fmt.Sprintf("%v", val))
	}

	return func() {
		for env, _ := range envMap {
			os.Unsetenv(env)
		}
	}
}

func stringPointer(s string) *string {
	return &s
}

func intPointer(i int) *int {
	return &i
}

func toJSONStr(object interface{}) string {
	data, _ := json.Marshal(object)
	return string(data)
}
