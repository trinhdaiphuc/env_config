package env_config

import (
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type optionalInner struct {
	Name string `env:"NAME"`
}

type optionalConfig struct {
	Section  *optionalInner `env:"SECTION"`
	Level    *string        `env:"LEVEL"`
	Retries  *int           `env:"RETRIES;default=3"`
	Required string         `env:"REQUIRED"`
}

func TestLoadConfigKeepsOptionalPointersNil(t *testing.T) {
	tests := []struct {
		name   string
		envs   map[string]string
		assert func(t *testing.T, cfg *optionalConfig)
	}{
		{
			name: "absent env leaves pointers nil",
			assert: func(t *testing.T, cfg *optionalConfig) {
				assert.Nil(t, cfg.Section)
				assert.Nil(t, cfg.Level)
			},
		},
		{
			name: "nested section allocated only when a child env is set",
			envs: map[string]string{"SECTION_NAME": "inner"},
			assert: func(t *testing.T, cfg *optionalConfig) {
				require.NotNil(t, cfg.Section)
				assert.Equal(t, "inner", cfg.Section.Name)
				assert.Nil(t, cfg.Level)
			},
		},
		{
			name: "leaf pointer allocated when its env is set",
			envs: map[string]string{"LEVEL": "debug"},
			assert: func(t *testing.T, cfg *optionalConfig) {
				require.NotNil(t, cfg.Level)
				assert.Equal(t, "debug", *cfg.Level)
				assert.Nil(t, cfg.Section)
			},
		},
		{
			name: "leaf pointer with a default is always allocated",
			assert: func(t *testing.T, cfg *optionalConfig) {
				require.NotNil(t, cfg.Retries)
				assert.Equal(t, 3, *cfg.Retries)
			},
		},
		{
			name: "empty env value still allocates",
			envs: map[string]string{"LEVEL": ""},
			assert: func(t *testing.T, cfg *optionalConfig) {
				require.NotNil(t, cfg.Level)
				assert.Equal(t, "", *cfg.Level)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			cfg := &optionalConfig{}
			require.NoError(t, LoadConfig(cfg))
			tt.assert(t, cfg)
		})
	}
}

func TestLoadConfigPreservesPreAllocatedSection(t *testing.T) {
	t.Setenv("SECTION_NAME", "inner")

	cfg := &optionalConfig{Section: &optionalInner{}}
	require.NoError(t, LoadConfig(cfg))
	assert.Equal(t, "inner", cfg.Section.Name)
}

type doublePointerConfig struct {
	Nested **optionalInner `env:"NESTED"`
}

type mapConfig struct {
	Values map[string]string `env:"MAPV"`
}

type durationSliceConfig struct {
	Timeouts []time.Duration `env:"TIMEOUTS"`
}

type int8Config struct {
	Value int8 `env:"V8"`
}

type intSliceConfig struct {
	Values []int `env:"VS"`
}

type boolSliceConfig struct {
	Values []bool `env:"BS"`
}

type delimiterOnScalarConfig struct {
	Value int `env:"VDI;delimiter=,"`
}

type unexportedConfig struct {
	name string `env:"NAME2"`
}

func TestLoadConfigReportsErrorsInsteadOfFailingSilently(t *testing.T) {
	tests := []struct {
		name    string
		cfg     interface{}
		envs    map[string]string
		wantErr string
	}{
		{
			name:    "pointer to pointer section",
			cfg:     &doublePointerConfig{},
			wantErr: "unsupported nested type",
		},
		{
			name:    "unsupported field type",
			cfg:     &mapConfig{},
			envs:    map[string]string{"MAPV": "a:b"},
			wantErr: "unsupported field type",
		},
		{
			name:    "unsupported slice element type",
			cfg:     &durationSliceConfig{},
			envs:    map[string]string{"TIMEOUTS": "1s,2s"},
			wantErr: "unsupported slice element type time.Duration",
		},
		{
			name:    "int8 overflow",
			cfg:     &int8Config{},
			envs:    map[string]string{"V8": "300"},
			wantErr: "value out of range",
		},
		{
			name:    "malformed int slice element",
			cfg:     &intSliceConfig{},
			envs:    map[string]string{"VS": "1,abc,3"},
			wantErr: `parsing element "abc"`,
		},
		{
			name:    "malformed bool slice element",
			cfg:     &boolSliceConfig{},
			envs:    map[string]string{"BS": "true,maybe"},
			wantErr: `parsing element "maybe"`,
		},
		{
			name:    "delimiter option on a scalar field",
			cfg:     &delimiterOnScalarConfig{},
			envs:    map[string]string{"VDI": "42"},
			wantErr: "delimiter option is not supported",
		},
		{
			name:    "unexported tagged field",
			cfg:     &unexportedConfig{},
			envs:    map[string]string{"NAME2": "hi"},
			wantErr: "cannot load unexported field",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for key, value := range tt.envs {
				t.Setenv(key, value)
			}

			err := LoadConfig(tt.cfg)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestLoadConfigNumericSliceDefaultWithoutDelimiter(t *testing.T) {
	type cfgType struct {
		Values []int `env:"NUMS;default=1,2,3"`
	}

	cfg := &cfgType{}
	require.NoError(t, LoadConfig(cfg))
	assert.Equal(t, []int{1, 2, 3}, cfg.Values)
}

type customStrategyType struct{}

type customStrategy struct{}

func (customStrategy) SetValue(reflect.Value, string, TagOption) error { return nil }

// TestRegisterStrategyIsRaceFree only fails under -race, which CI runs.
func TestRegisterStrategyIsRaceFree(t *testing.T) {
	t.Setenv("LEVEL", "debug")

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			RegisterStrategy(reflect.TypeOf(customStrategyType{}), customStrategy{})
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			if err := LoadConfig(&optionalConfig{}); err != nil {
				t.Error(err)
				return
			}
		}
	}()

	wg.Wait()
}
