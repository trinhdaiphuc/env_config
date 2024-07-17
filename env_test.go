package env_config

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"testing"
	"time"
)

type Config struct {
	Scheduler                  *Scheduler                  `env:"SCHEDULER"`
	ServerConfig               *ServerConfig               `env:"SERVER"`
	NestedConfig               *NestedConfig               `env:"APP"`
	UninitializedPointerConfig *UninitializedPointerConfig `env:"UN_INIT"`
	EnvVarTaggedConfig         *EnvVarTaggedConfig         `env:"TAG"`
	Logger                     *Logger                     `env:"LOGGER"`
	UnFieldConfig              string
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
	AppName string `env:"APP_NAME"`
	Debug   bool   `env:"DEBUG"`
}

type Logger struct {
	Level   string `env:"LEVEL"`
	Encoder string `env:"ENCODER"`
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
				},
				Logger: &Logger{
					Level:   "debug",
					Encoder: "json",
				},
				UnFieldConfig: "",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.args.cfg, tt.wantCfg) {
				t.Errorf("LoadConfig() \ngot  = %v \nwant = %v", toJSONStr(tt.args.cfg), toJSONStr(tt.wantCfg))
			}
		})
	}
}

func setEnv(timeNow time.Time) func() {
	envMap := map[string]interface{}{
		"SCHEDULER_INTERVAL":    "5m",
		"SCHEDULER_START_AT":    timeNow.Format(time.RFC3339),
		"SERVER_REDIS_HOST":     "127.0.0.1",
		"SERVER_REDIS_PORT":     6379,
		"SERVER_REDIS_PASSWORD": "secret",
		"APP_DB_HOST":           "localhost",
		"APP_DB_PORT":           3306,
		"APP_DB_USER":           "root",
		"APP_DB_PASS":           "password",
		"UN_INIT_LOG_LEVEL":     "info",
		"UN_INIT_TIMEOUT":       10,
		"TAG_APP_NAME":          "env_config",
		"TAG_DEBUG":             true,
		"LOGGER_LEVEL":          "debug",
		"LOGGER_ENCODER":        "json",
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
