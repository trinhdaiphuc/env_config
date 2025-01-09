package env_config

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

type tagOptionError struct{}

func (t tagOptionError) Next() TagOption {
	return nil
}

func (t tagOptionError) Apply(value string) (interface{}, error) {
	return nil, errors.New("error")
}

func (t tagOptionError) SetValue(value string) {
}

func (t tagOptionError) SetNext(option TagOption) {
}

var (
	_ TagOption = tagOptionError{}
)

func TestStringStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic string",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "test_value",
			},
			want:    reflect.ValueOf("test_value"),
			wantErr: false,
		},
		{
			name: "empty string",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf(""),
			wantErr: false,
		},
		{
			name: "invalid string",
			args: args{
				field:    reflect.New(reflect.TypeOf(func() {})).Elem(),
				envValue: "",
			},
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf("")).Elem(),
				envValue:  "test_value",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestIntStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "positive integer",
			args: args{
				field:    reflect.New(reflect.TypeOf(int64(0))).Elem(),
				envValue: "123",
			},
			want:    reflect.ValueOf(int64(123)),
			wantErr: false,
		},
		{
			name: "negative integer",
			args: args{
				field:    reflect.New(reflect.TypeOf(int64(0))).Elem(),
				envValue: "-456",
			},
			want:    reflect.ValueOf(int64(-456)),
			wantErr: false,
		},
		{
			name: "invalid integer",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(int64(0)),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(int64(0))).Elem(),
				envValue:  "123",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IntStrategy[int64]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestUintStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "positive unsigned integer",
			args: args{
				field:    reflect.New(reflect.TypeOf(uint64(0))).Elem(),
				envValue: "123",
			},
			want:    reflect.ValueOf(uint64(123)),
			wantErr: false,
		},
		{
			name: "invalid unsigned integer",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(uint64(0)),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(uint64(0))).Elem(),
				envValue:  "123",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UintStrategy[uint64]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestFloatStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "positive float",
			args: args{
				field:    reflect.New(reflect.TypeOf(0.0)).Elem(),
				envValue: "123.45",
			},
			want:    reflect.ValueOf(123.45),
			wantErr: false,
		},
		{
			name: "negative float",
			args: args{
				field:    reflect.New(reflect.TypeOf(0.0)).Elem(),
				envValue: "-678.90",
			},
			want:    reflect.ValueOf(-678.90),
			wantErr: false,
		},
		{
			name: "invalid float",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(0.0),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(0.0)).Elem(),
				envValue:  "123.45",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FloatStrategy[float64]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestBoolStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "true",
			args: args{
				field:    reflect.New(reflect.TypeOf(false)).Elem(),
				envValue: "true",
			},
			want:    reflect.ValueOf(true),
			wantErr: false,
		},
		{
			name: "false",
			args: args{
				field:    reflect.New(reflect.TypeOf(true)).Elem(),
				envValue: "false",
			},
			want:    reflect.ValueOf(false),
			wantErr: false,
		},
		{
			name: "invalid bool",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(false),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(false)).Elem(),
				envValue:  "true",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BoolStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestByteSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic byte slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]byte{})).Elem(),
				envValue: "test_value",
			},
			want:    reflect.ValueOf([]byte("test_value")),
			wantErr: false,
		},
		{
			name: "empty byte slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]byte{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]byte("")),
			wantErr: false,
		},
		{
			name: "invalid byte slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]byte("")),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]byte{})).Elem(),
				envValue:  "test_value",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ByteSliceStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestDurationStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic duration",
			args: args{
				field:    reflect.New(reflect.TypeOf(time.Duration(0))).Elem(),
				envValue: "1s",
			},
			want:    reflect.ValueOf(time.Second),
			wantErr: false,
		},
		{
			name: "empty duration",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf(0),
			wantErr: false,
		},
		{
			name: "invalid duration",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(0),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(0)).Elem(),
				envValue:  "1s",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DurationStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestTimeStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic time",
			args: args{
				field:    reflect.New(reflect.TypeOf(time.Time{})).Elem(),
				envValue: "2021-01-02T15:04:05Z",
			},
			want:    reflect.ValueOf(time.Date(2021, 1, 2, 15, 4, 5, 0, time.UTC)),
			wantErr: false,
		},
		{
			name: "empty time",
			args: args{
				field:    reflect.New(reflect.TypeOf(time.Time{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf(time.Time{}),
			wantErr: false,
		},
		{
			name: "invalid time",
			args: args{
				field:    reflect.New(reflect.TypeOf(0)).Elem(),
				envValue: "invalid",
			},
			want:    reflect.ValueOf(time.Time{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf(time.Time{})).Elem(),
				envValue:  "2021-01-02T15:04:05Z",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := TimeStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func Test_parseOptionValue(t *testing.T) {
	type args struct {
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "basic",
			args: args{
				envValue: "test_value",
			},
			want:    "test_value",
			wantErr: false,
		},
		{
			name: "parse option value with DefaultOption",
			args: args{
				envValue:  "",
				tagOption: &DefaultOption{DefaultValue: "default_value"},
			},
			want:    "default_value",
			wantErr: false,
		},
		{
			name: "error tag option",
			args: args{
				envValue:  "test_value",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseOptionValue(tt.args.envValue, tt.args.tagOption)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseOptionValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("parseOptionValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	var (
		nilStringSlice []string
	)
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic string slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]string{})).Elem(),
				envValue: "test_value1,test_value2",
			},
			want:    reflect.ValueOf([]string{"test_value1", "test_value2"}),
			wantErr: false,
		},
		{
			name: "empty string slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]string{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf(nilStringSlice),
			wantErr: false,
		},
		{
			name: "invalid string slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]string{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]string{})).Elem(),
				envValue:  "test_value1,test_value2",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := StringSliceStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestBoolSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic bool slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]bool{})).Elem(),
				envValue: "true,false",
			},
			want:    reflect.ValueOf([]bool{true, false}),
			wantErr: false,
		},
		{
			name: "empty bool slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]bool{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]bool{}),
			wantErr: false,
		},
		{
			name: "invalid bool slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]bool{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]bool{})).Elem(),
				envValue:  "true,false",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BoolSliceStrategy{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestIntSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic int slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]int{})).Elem(),
				envValue: "123,-456",
			},
			want:    reflect.ValueOf([]int{123, -456}),
			wantErr: false,
		},
		{
			name: "empty int slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]int{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]int{}),
			wantErr: false,
		},
		{
			name: "invalid int slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]int{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]int{})).Elem(),
				envValue:  "123,-456",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := IntSliceStrategy[int]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestUintSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic uint slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]uint{})).Elem(),
				envValue: "123,456",
			},
			want:    reflect.ValueOf([]uint{123, 456}),
			wantErr: false,
		},
		{
			name: "empty uint slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]uint{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]uint{}),
			wantErr: false,
		},
		{
			name: "invalid uint slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]uint{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]uint{})).Elem(),
				envValue:  "123,456",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := UintSliceStrategy[uint]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

func TestFloatSliceStrategy_SetValue(t *testing.T) {
	type args struct {
		field     reflect.Value
		envValue  string
		tagOption TagOption
	}
	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "basic float slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]float64{})).Elem(),
				envValue: "123.45,-678.90",
			},
			want:    reflect.ValueOf([]float64{123.45, -678.90}),
			wantErr: false,
		},
		{
			name: "empty float slice",
			args: args{
				field:    reflect.New(reflect.TypeOf([]float64{})).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]float64{}),
			wantErr: false,
		},
		{
			name: "invalid float slice",
			args: args{
				field:    reflect.New(reflect.TypeOf("")).Elem(),
				envValue: "",
			},
			want:    reflect.ValueOf([]float64{}),
			wantErr: true,
		},
		{
			name: "error tag option",
			args: args{
				field:     reflect.New(reflect.TypeOf([]float64{})).Elem(),
				envValue:  "123.45,-678.90",
				tagOption: tagOptionError{},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := FloatSliceStrategy[float64]{}
			if err := s.SetValue(tt.args.field, tt.args.envValue, tt.args.tagOption); (err != nil) != tt.wantErr {
				t.Errorf("SetValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && !reflect.DeepEqual(tt.args.field.Interface(), tt.want.Interface()) {
				t.Errorf("value got %v, want %v", tt.args.field.Interface(), tt.want.Interface())
			}
		})
	}
}

// Define a custom type and strategy for testing
type customType struct {
	Value string
}

type customTypeStrategy struct{}

func (s customTypeStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	field.Set(reflect.ValueOf(customType{Value: envValue}))
	return nil
}
func TestRegisterStrategy(t *testing.T) {
	type args struct {
		strategyType reflect.Type
		strategy     TypeStrategy
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantRegistered bool
	}{
		{
			name: "Register custom strategy",
			args: args{
				strategyType: reflect.TypeOf(customType{}),
				strategy:     customTypeStrategy{},
			},
			wantErr:        false,
			wantRegistered: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RegisterStrategy(tt.args.strategyType, tt.args.strategy)

			registeredStrategy, exists := complexTypeStrategies[tt.args.strategyType]
			if exists != tt.wantRegistered {
				t.Errorf("expected registered = %v, got %v", tt.wantRegistered, exists)
			}

			if tt.wantRegistered {
				if !reflect.DeepEqual(registeredStrategy, tt.args.strategy) {
					t.Errorf("expected strategy = %v, got %v", tt.args.strategy, registeredStrategy)
				}
			}
		})
	}
}
