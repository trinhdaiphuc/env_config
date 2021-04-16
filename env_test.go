package env_config

import (
	"os"
	"testing"
)

func TestHello(t *testing.T) {
	os.Setenv("HELLO", "")
	if hello:=Env("HELLO", "world"); hello != "world" {
		t.Errorf("Output expect HELLO variable value world instead of %v", hello)
	}

	os.Setenv("HI", "hello")

	if hi:=Env("HELLO", "world"); hi != "world" {
		t.Errorf("Output expect HI variable value world instead of %v", hi)
	}
}
