package env_config

import (
	"os"
)

func LoadConfig(cfg interface{}) error {
	root := NewStruct(cfg, "")
	return root.Load()
}

func Env(key, defaultValue string) (value string) {
	if value = os.Getenv(key); value == "" {
		value = defaultValue
	}
	return
}
