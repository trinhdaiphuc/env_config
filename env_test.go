package env_config

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	type args struct {
		cfg interface{}
	}
	var (
		host     = "127.0.0.1"
		port     = "6379"
		password = "secret"
	)
	os.Setenv("REDIS_HOST", host)
	os.Setenv("REDIS_PORT", port)
	os.Setenv("REDIS_PASSWORD", password)
	defer os.Unsetenv("REDIS_HOST")
	defer os.Unsetenv("REDIS_PORT")
	defer os.Unsetenv("REDIS_PASSWORD")

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test LoadConfig",
			args: args{
				cfg: &ServerConfig{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := LoadConfig(tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			data, _ := json.Marshal(tt.args.cfg)
			fmt.Println("Config ", string(data))
		})
	}
}
