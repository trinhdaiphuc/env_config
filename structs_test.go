package env_config

import (
	"os"
	"reflect"
	"testing"
)

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
		{
			name:         "Test empty tag",
			tag:          "",
			expectedKey:  "",
			expectedTags: nil,
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
				envVarKey:   "REDIS_HOST",
				envVarValue: "localhost",
			},
			want:    "localhost",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv(tt.args.envVarKey, tt.args.envVarValue)
			defer os.Unsetenv(tt.args.envVarKey)

			structItem, err := NewStruct(tt.args.cfg, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("StructItem.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			err = structItem.Load()
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
			name: "Test FieldItem Load with RedisConfig",
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

			structItem, err := NewStruct(tt.args.cfg, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("StructItem.Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			for _, child := range structItem.Children() {
				if configItem, ok := child.(FieldItem); ok && configItem.Key() == tt.args.configKey {
					err := configItem.Load()
					if (err != nil) != tt.wantErr {
						t.Errorf("FieldItem.Load() error = %v, wantErr %v", err, tt.wantErr)
					}

					if configItem.Value().String() != tt.want {
						t.Errorf("Expected %s, got %s", tt.want, configItem.Value().String())
					}
				}
			}
		})
	}
}

func TestNewStruct(t *testing.T) {
	type args struct {
		s         interface{}
		keyPrefix string
	}
	tests := []struct {
		name    string
		args    args
		want    StructItem
		wantErr bool
	}{
		{
			name: "Valid ServerConfig without prefix",
			args: args{
				s:         &ServerConfig{},
				keyPrefix: "",
			},
			want: StructItem{
				prefix: "",
				children: []Item{
					StructItem{
						prefix:     "REDIS",
						raw:        &RedisConfig{},
						value:      reflect.ValueOf(&RedisConfig{}),
						tagOptions: nil,
						children: []Item{
							FieldItem{
								key:        "REDIS_HOST",
								tagOptions: nil,
							},
							FieldItem{
								key:        "REDIS_PORT",
								tagOptions: nil,
							},
							FieldItem{
								key:        "REDIS_PASSWORD",
								tagOptions: nil,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Valid NestedConfig with prefix",
			args: args{
				s:         &NestedConfig{},
				keyPrefix: "APP_",
			},
			want: StructItem{
				prefix: "APP_",
				value:  reflect.ValueOf(&NestedConfig{}),
				children: []Item{
					StructItem{
						prefix: "APP_DB",
						value:  reflect.ValueOf(&DatabaseConfig{}),
						children: []Item{
							FieldItem{
								raw:        "",
								key:        "APP_DB_HOST",
								value:      reflect.ValueOf(""),
								fieldName:  "Host",
								tagOptions: nil,
							},
							FieldItem{
								raw:        0,
								key:        "APP_DB_PORT",
								value:      reflect.ValueOf(0),
								fieldName:  "Port",
								tagOptions: nil,
							},
							FieldItem{
								raw:        "",
								key:        "APP_DB_USER",
								value:      reflect.ValueOf(""),
								fieldName:  "Username",
								tagOptions: nil,
							},
							FieldItem{
								raw:        "",
								key:        "APP_DB_PASS",
								value:      reflect.ValueOf(""),
								fieldName:  "Password",
								tagOptions: nil,
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "UninitializedPointerConfig without prefix",
			args: args{
				s:         &UninitializedPointerConfig{},
				keyPrefix: "",
			},
			want: StructItem{
				prefix: "",
				children: []Item{
					FieldItem{
						key:        "LOG_LEVEL",
						tagOptions: nil,
					},
					FieldItem{
						key:        "TIMEOUT",
						tagOptions: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "EnvVarTaggedConfig with prefix",
			args: args{
				s:         &EnvVarTaggedConfig{},
				keyPrefix: "APP_",
			},
			want: StructItem{
				prefix: "APP_",
				children: []Item{
					FieldItem{
						key:        "APP_APP_NAME",
						tagOptions: nil,
					},
					FieldItem{
						key:        "APP_DEBUG",
						tagOptions: nil,
					},
					FieldItem{
						key:        "APP_PI",
						tagOptions: nil,
					},
					FieldItem{
						key:        "APP_NUMBER",
						tagOptions: nil,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Non-struct input",
			args: args{
				s:         123,
				keyPrefix: "",
			},
			want:    StructItem{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewStruct(tt.args.s, tt.args.keyPrefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !compareStructItems(got, tt.want) {
				t.Errorf("NewStruct() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func compareStructItems(a, b StructItem) bool {
	if a.prefix != b.prefix || !reflect.DeepEqual(a.tagOptions, b.tagOptions) || len(a.children) != len(b.children) {
		return false
	}
	for i, childA := range a.children {
		childB := b.children[i]
		structChildA, ok := childA.(StructItem)
		if !ok {
			if !compareItems(childA, childB) {
				return false
			}
		} else {
			structChildB := childB.(StructItem)
			if !compareStructItems(structChildA, structChildB) {
				return false
			}
		}
	}
	return true
}

func compareItems(a, b Item) bool {
	if a.Key() != b.Key() || !reflect.DeepEqual(a.TagOptions(), b.TagOptions()) {
		return false
	}
	return true
}
