package env_config

import (
	"reflect"
	"strconv"
	"time"
)

// TypeStrategy is an interface that defines a method to set a value to a field
type TypeStrategy interface {
	SetValue(field reflect.Value, envValue string) error
}

type StringStrategy struct{}

func (s StringStrategy) SetValue(field reflect.Value, envValue string) error {
	field.SetString(envValue)
	return nil
}

type IntStrategy struct{}

func (s IntStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := strconv.ParseInt(envValue, 10, 64)
	if err != nil {
		return err
	}
	field.SetInt(v)
	return nil
}

type UintStrategy struct{}

func (s UintStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := strconv.ParseUint(envValue, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}

type FloatStrategy struct{}

func (s FloatStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := strconv.ParseFloat(envValue, 64)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return nil
}

type BoolStrategy struct{}

func (s BoolStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := strconv.ParseBool(envValue)
	if err != nil {
		return err
	}
	field.SetBool(v)
	return nil
}

type ByteSliceStrategy struct{}

func (s ByteSliceStrategy) SetValue(field reflect.Value, envValue string) error {
	field.SetBytes([]byte(envValue))
	return nil
}

type DurationStrategy struct{}

func (s DurationStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := time.ParseDuration(envValue)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

type TimeStrategy struct{}

func (s TimeStrategy) SetValue(field reflect.Value, envValue string) error {
	v, err := time.Parse(time.RFC3339, envValue)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

var strategies = map[reflect.Kind]TypeStrategy{
	reflect.String:  StringStrategy{},
	reflect.Int:     IntStrategy{},
	reflect.Int8:    IntStrategy{},
	reflect.Int16:   IntStrategy{},
	reflect.Int32:   IntStrategy{},
	reflect.Int64:   IntStrategy{},
	reflect.Uint:    UintStrategy{},
	reflect.Uint8:   UintStrategy{},
	reflect.Uint16:  UintStrategy{},
	reflect.Uint32:  UintStrategy{},
	reflect.Uint64:  UintStrategy{},
	reflect.Float32: FloatStrategy{},
	reflect.Float64: FloatStrategy{},
	reflect.Bool:    BoolStrategy{},
	reflect.Slice:   ByteSliceStrategy{},
}
