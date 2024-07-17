package env_config

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
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
	if envValue == "" {
		return nil
	}

	v, err := strconv.ParseInt(envValue, 10, 64)
	if err != nil {
		return err
	}

	field.SetInt(v)
	return nil
}

type UintStrategy struct{}

func (s UintStrategy) SetValue(field reflect.Value, envValue string) error {
	if envValue == "" {
		return nil
	}

	v, err := strconv.ParseUint(envValue, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}

type FloatStrategy struct{}

func (s FloatStrategy) SetValue(field reflect.Value, envValue string) error {
	if envValue == "" {
		return nil
	}

	v, err := strconv.ParseFloat(envValue, 64)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return nil
}

type BoolStrategy struct{}

func (s BoolStrategy) SetValue(field reflect.Value, envValue string) error {
	if envValue == "" {
		return nil
	}
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
	if envValue == "" {
		return nil
	}
	v, err := time.ParseDuration(envValue)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

type TimeStrategy struct{}

func (s TimeStrategy) SetValue(field reflect.Value, envValue string) error {
	if envValue == "" {
		return nil
	}
	v, err := time.Parse(time.RFC3339, envValue)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

type StringSliceStrategy struct{}

func (s StringSliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("invalid type, expected []string but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	v.Set(reflect.ValueOf(values))
	return nil
}

type BoolSliceStrategy struct{}

func (s BoolSliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Bool {
		return fmt.Errorf("invalid type, expected []bool but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	boolValues := make([]bool, len(values))
	for i, val := range values {
		var b bool
		b, err := strconv.ParseBool(val)
		if err == nil {
			boolValues[i] = b
		}
	}
	v.Set(reflect.ValueOf(boolValues))
	return nil
}

type IntSliceStrategy struct{}

func (s IntSliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Int {
		return fmt.Errorf("invalid type, expected []int but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	intValues := make([]int, len(values))
	for i, val := range values {
		var b int
		b, err := strconv.Atoi(val)
		if err == nil {
			intValues[i] = b
		}
	}
	v.Set(reflect.ValueOf(intValues))
	return nil
}

type UintSliceStrategy struct{}

func (s UintSliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint {
		return fmt.Errorf("invalid type, expected []uint but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	uintValues := make([]uint, len(values))
	for i, val := range values {
		var b uint64
		b, err := strconv.ParseUint(val, 10, 64)
		if err == nil {
			uintValues[i] = uint(b)
		}
	}
	v.Set(reflect.ValueOf(uintValues))
	return nil
}

type Float64SliceStrategy struct{}

func (s Float64SliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Float64 {
		return fmt.Errorf("invalid type, expected []float64 but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	floatValues := make([]float64, len(values))
	for i, val := range values {
		var b float64
		b, err := strconv.ParseFloat(val, 64)
		if err == nil {
			floatValues[i] = b
		}
	}
	v.Set(reflect.ValueOf(floatValues))
	return nil
}

type Float32SliceStrategy struct{}

func (s Float32SliceStrategy) SetValue(v reflect.Value, value string) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Float32 {
		return fmt.Errorf("invalid type, expected []float32 but got %s", v.Kind())
	}

	values := strings.Split(value, ",")
	floatValues := make([]float32, len(values))
	for i, val := range values {
		floatValue, err := strconv.ParseFloat(val, 64)
		if err == nil {
			floatValues[i] = float32(floatValue)
		}
	}
	v.Set(reflect.ValueOf(floatValues))
	return nil
}

var buildInTypeStrategies = map[reflect.Kind]TypeStrategy{
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

var complexTypeStrategies = map[reflect.Type]TypeStrategy{
	reflect.TypeOf(time.Duration(0)): DurationStrategy{},
	reflect.TypeOf(time.Time{}):      TimeStrategy{},
	reflect.TypeOf([]string{}):       StringSliceStrategy{},
	reflect.TypeOf([]bool{}):         BoolSliceStrategy{},
	reflect.TypeOf([]int{}):          IntSliceStrategy{},
	reflect.TypeOf([]uint{}):         UintSliceStrategy{},
	reflect.TypeOf([]float64{}):      Float64SliceStrategy{},
	reflect.TypeOf([]float32{}):      Float32SliceStrategy{},
}
