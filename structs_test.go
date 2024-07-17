package env_config

import (
	"testing"
)

type ConfigStruct struct {
	Host string `env:"HOST,localhost"`
	Port int    `env:"PORT,8080"`
}

func TestStruct(t *testing.T) {
	var (
		cfg = ConfigStruct{}
		out = Map(cfg)
	)
	t.Logf("out %v", out)
}
