package env_config

import (
	"os"
	"strconv"
	"testing"
	"time"
)

type Config struct {
	Host    string        `env:"HOST,localhost"`
	Port    int           `env:"PORT,8080"`
	Timeout time.Duration `env:"TIMEOUT"`
}

func TestHello(t *testing.T) {
	os.Setenv("HELLO", "")
	if hello := Env("HELLO", "world"); hello != "world" {
		t.Errorf("Output expect HELLO variable value world instead of %v", hello)
	}

	os.Setenv("HI", "hello")

	if hi := Env("HELLO", "world"); hi != "world" {
		t.Errorf("Output expect HI variable value world instead of %v", hi)
	}
}

func TestEnvStruct(t *testing.T) {
	var (
		host         = "127.0.0.1"
		port         = "8081"
		timeout      = "1h30m"
		portInt, _   = strconv.ParseInt(port, 10, 64)
		timeParse, _ = time.ParseDuration(timeout)
	)
	os.Setenv("HOST", host)
	os.Setenv("PORT", port)
	os.Setenv("TIMEOUT", timeout)
	cfg := &Config{}
	err := EnvStruct(cfg)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Config %+v", cfg)
	if cfg.Host != host {
		t.Errorf("Output expect %v variable value world instead of %v", host, cfg.Host)
	}
	if cfg.Port != int(portInt) {
		t.Errorf("Output expect %v variable value world instead of %d", port, cfg.Port)
	}
	if cfg.Timeout != timeParse {
		t.Errorf("Output expect %v variable value world instead of %d", timeout, cfg.Timeout)
	}

	os.Clearenv()

	err = EnvStruct(cfg)
	if err != nil {
		t.Error(err)
	}
	t.Logf("Config %+v", cfg)

	if cfg.Host != "localhost" {
		t.Errorf("Output expect localhost variable value world instead of %v", cfg.Host)
	}
	if cfg.Port != 8080 {
		t.Errorf("Output expect 8080 variable value world instead of %d", cfg.Port)
	}
}
