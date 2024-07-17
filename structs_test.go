package env_config

import (
	"os"
	"reflect"
	"testing"
)

type RedisConfig struct {
	Host     string `env:"HOST"`
	Port     int    `env:"PORT"`
	Password string `env:"PASSWORD"`
}

type ServerConfig struct {
	CacheConfig *RedisConfig `env:"REDIS"`
}

func TestNewStruct(t *testing.T) {
	tests := []struct {
		name      string
		cfg       interface{}
		keyPrefix string
		wantKind  reflect.Kind
		wantErr   bool
	}{
		{
			name:      "Test NewStruct with ServerConfig",
			cfg:       &ServerConfig{},
			keyPrefix: "",
			wantKind:  reflect.Struct,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			structItem := NewStruct(tt.cfg, tt.keyPrefix)
			if structItem.Value().Kind() != tt.wantKind {
				t.Errorf("Expected StructItem to hold a struct value, got %s", structItem.Value().Kind())
			}
		})
	}
}

func TestParseTagAndKey(t *testing.T) {
	tests := []struct {
		name         string
		tag          string
		expectedKey  string
		expectedTags []TagOption
	}{
		{
			name:         "Test ParseTagAndKey with single tag",
			tag:          "CACHE_REDIS_HOST,key=value",
			expectedKey:  "CACHE_REDIS_HOST",
			expectedTags: []TagOption{{key: "key", value: "value"}},
		},
		{
			name:         "Test ParseTagAndKey with multiple tags",
			tag:          "CACHE_REDIS_HOST,key1=value1,key2=value2",
			expectedKey:  "CACHE_REDIS_HOST",
			expectedTags: []TagOption{{key: "key1", value: "value1"}, {key: "key2", value: "value2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key, tagOptions := parseTagAndKey(tt.tag)
			if key != tt.expectedKey {
				t.Errorf("Expected key %s, got %s", tt.expectedKey, key)
			}

			if !reflect.DeepEqual(tagOptions, tt.expectedTags) {
				t.Errorf("Expected tag options %v, got %v", tt.expectedTags, tagOptions)
			}
		})
	}
}

func TestStructItem_Load(t *testing.T) {
	type args struct {
		cfg         interface{}
		envVarKey   string
		envVarValue string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test StructItem Load with ServerConfig",
			args: args{
				cfg:         &ServerConfig{},
				envVarKey:   "CACHE_REDIS_HOST",
				envVarValue: "cache_redis_host_value",
			},
			want:    "cache_redis_host_value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.envVarKey, tt.args.envVarValue)
			defer os.Unsetenv(tt.args.envVarKey)

			structItem := NewStruct(tt.args.cfg, "")
			err := structItem.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("StructItem.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			serverConfig := tt.args.cfg.(*ServerConfig)
			if serverConfig.CacheConfig.Host != tt.want {
				t.Errorf("Expected %s, got %s", tt.want, serverConfig.CacheConfig.Host)
			}
		})
	}
}

func TestCombineKeyPrefix(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		key      string
		expected string
	}{
		{
			name:     "Test CombineKeyPrefix with empty prefix",
			prefix:   "",
			key:      "KEY",
			expected: "KEY",
		},
		{
			name:     "Test CombineKeyPrefix with prefix ending with underscore",
			prefix:   "PREFIX_",
			key:      "KEY",
			expected: "PREFIX_KEY",
		},
		{
			name:     "Test CombineKeyPrefix with prefix not ending with underscore",
			prefix:   "PREFIX",
			key:      "KEY",
			expected: "PREFIX_KEY",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := combineKeyPrefix(tt.prefix, tt.key)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestConfigItem_Load(t *testing.T) {
	type args struct {
		cfg         interface{}
		envVarKey   string
		envVarValue string
		configKey   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Test ConfigItem Load with RedisConfig",
			args: args{
				cfg:         &RedisConfig{},
				envVarKey:   "REDIS_HOST",
				envVarValue: "redis_host_value",
				configKey:   "REDIS_HOST",
			},
			want:    "redis_host_value",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.envVarKey, tt.args.envVarValue)
			defer os.Unsetenv(tt.args.envVarKey)

			structItem := NewStruct(tt.args.cfg, "")
			for _, child := range structItem.Children() {
				if configItem, ok := child.(ConfigItem); ok && configItem.Key() == tt.args.configKey {
					err := configItem.Load()
					if (err != nil) != tt.wantErr {
						t.Errorf("ConfigItem.Load() error = %v, wantErr %v", err, tt.wantErr)
					}

					if configItem.Value().String() != tt.want {
						t.Errorf("Expected %s, got %s", tt.want, configItem.Value().String())
					}
				}
			}
		})
	}
}
