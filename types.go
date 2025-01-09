package env_config

import (
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// TypeStrategy is an interface that defines a method to set a value to a field
type TypeStrategy interface {
	SetValue(field reflect.Value, envValue string, tagOption TagOption) error
}

func RegisterStrategy(strategyType reflect.Type, strategy TypeStrategy) {
	complexTypeStrategies[strategyType] = strategy
}

var (
	complexTypeStrategies = make(map[reflect.Type]TypeStrategy)
	buildInTypeStrategies = make(map[reflect.Kind]TypeStrategy)
)

type StringStrategy struct{}

func (s StringStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.String {
		return fmt.Errorf("invalid type, expected string but got %s", field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	field.SetString(value)
	return nil
}

type IntStrategy[I IntType] struct{}

func (s IntStrategy[I]) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.TypeFor[I]().Kind() {
		return fmt.Errorf("invalid type, expected %v but got %s", reflect.TypeFor[I]().Kind(), field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	v, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}

	field.SetInt(v)
	return nil
}

type UintStrategy[U UintType] struct{}

func (s UintStrategy[U]) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.TypeFor[U]().Kind() {
		return fmt.Errorf("invalid type, expected %v but got %s", reflect.TypeFor[U]().Kind(), field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	v, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return err
	}
	field.SetUint(v)
	return nil
}

type FloatStrategy[F FloatType] struct{}

func (s FloatStrategy[F]) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.TypeFor[F]().Kind() {
		return fmt.Errorf("invalid type, expected %v but got %s", reflect.TypeFor[F]().Kind(), field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}

	v, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	field.SetFloat(v)
	return nil
}

type BoolStrategy struct{}

func (s BoolStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.Bool {
		return fmt.Errorf("invalid type, expected bool but got %s", field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}
	v, err := strconv.ParseBool(value)
	if err != nil {
		return err
	}
	field.SetBool(v)
	return nil
}

type ByteSliceStrategy struct{}

func (s ByteSliceStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.Slice || field.Type().Elem().Kind() != reflect.Uint8 {
		return fmt.Errorf("invalid type, expected []byte but got %s", field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	field.SetBytes([]byte(value))
	return nil
}

type DurationStrategy struct{}

func (s DurationStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}
	v, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

type TimeStrategy struct{}

func (s TimeStrategy) SetValue(field reflect.Value, envValue string, tagOption TagOption) error {
	if field.Kind() != reflect.Struct || field.Type() != reflect.TypeOf(time.Time{}) {
		return fmt.Errorf("invalid type, expected time.Time but got %s", field.Kind())
	}

	value, err := parseOptionValue(envValue, tagOption)
	if err != nil {
		return err
	}

	if value == "" {
		return nil
	}
	v, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return err
	}
	field.Set(reflect.ValueOf(v))
	return nil
}

type StringSliceStrategy struct{}

func (s StringSliceStrategy) SetValue(v reflect.Value, envValue string, tagOption TagOption) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.String {
		return fmt.Errorf("invalid type, expected []string but got %s", v.Kind())
	}

	tagOption = setStringSliceDefaultTagOption(tagOption)
	values, err := parseOptionValues(envValue, tagOption)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(values))
	return nil
}

func setStringSliceDefaultTagOption(tagOption TagOption) TagOption {
	if tagOption == nil {
		return tagOption
	}

	tag := tagOption
	hasDelimiter := false
	for tag != nil {
		if _, ok := tag.(*DelimiterOption); ok {
			hasDelimiter = true
			break
		}
		tag = tag.Next()
	}

	if !hasDelimiter {
		tagOption.SetNext(defaultTagOption())
	}

	return tagOption
}

type BoolSliceStrategy struct{}

func (s BoolSliceStrategy) SetValue(v reflect.Value, envValue string, tagOption TagOption) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Bool {
		return fmt.Errorf("invalid type, expected []bool but got %s", v.Kind())
	}

	values, err := parseOptionValues(envValue, tagOption)
	if err != nil {
		return err
	}

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

type IntSliceStrategy[I IntType] struct{}

func (s IntSliceStrategy[I]) SetValue(v reflect.Value, envValue string, tagOption TagOption) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.TypeFor[I]().Kind() {
		return fmt.Errorf("invalid type, expected []%v but got %s", reflect.TypeFor[I]().Kind(), v.Kind())
	}

	values, err := parseOptionValues(envValue, tagOption)
	if err != nil {
		return err
	}

	intValues := StringArrayToIntArray[I](values)
	v.Set(reflect.ValueOf(intValues))
	return nil
}

type UintSliceStrategy[U UintType] struct{}

func (s UintSliceStrategy[U]) SetValue(v reflect.Value, envValue string, tagOption TagOption) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.TypeFor[U]().Kind() {
		return fmt.Errorf("invalid type, expected []%v but got %s", reflect.TypeFor[U]().Kind(), v.Kind())
	}

	values, err := parseOptionValues(envValue, tagOption)
	if err != nil {
		return err
	}

	uintValues := StringArrayToUintArray[U](values)
	v.Set(reflect.ValueOf(uintValues))
	return nil
}

type FloatSliceStrategy[F FloatType] struct{}

func (s FloatSliceStrategy[F]) SetValue(v reflect.Value, envValue string, tagOption TagOption) error {
	if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.TypeFor[F]().Kind() {
		return fmt.Errorf("invalid type, expected []%v but got %s", reflect.TypeFor[F]().Kind(), v.Kind())
	}

	values, err := parseOptionValues(envValue, tagOption)
	if err != nil {
		return err
	}

	floatValues := StringArrayToFloatArray[F](values)
	v.Set(reflect.ValueOf(floatValues))
	return nil
}

func parseOptionValue(envValue string, option TagOption) (string, error) {
	if option == nil {
		return envValue, nil
	}
	value, err := option.Apply(envValue)
	if err != nil {
		return "", err
	}

	envValue, _ = value.(string)
	return envValue, nil
}

func parseOptionValues(envValue string, option TagOption) ([]string, error) {
	if option == nil {
		return parseOptionValues(envValue, defaultTagOption())
	}

	value, err := option.Apply(envValue)
	if err != nil {
		return nil, err
	}

	valueArr, _ := value.([]string)
	return valueArr, nil
}
func init() {
	buildInTypeStrategies = map[reflect.Kind]TypeStrategy{
		reflect.String:  StringStrategy{},
		reflect.Int:     IntStrategy[int]{},
		reflect.Int8:    IntStrategy[int8]{},
		reflect.Int16:   IntStrategy[int16]{},
		reflect.Int32:   IntStrategy[int32]{},
		reflect.Int64:   IntStrategy[int64]{},
		reflect.Uint:    UintStrategy[uint]{},
		reflect.Uint8:   UintStrategy[uint8]{},
		reflect.Uint16:  UintStrategy[uint16]{},
		reflect.Uint32:  UintStrategy[uint32]{},
		reflect.Uint64:  UintStrategy[uint64]{},
		reflect.Float32: FloatStrategy[float32]{},
		reflect.Float64: FloatStrategy[float64]{},
		reflect.Bool:    BoolStrategy{},
		reflect.Slice:   ByteSliceStrategy{},
	}

	complexTypeStrategies = map[reflect.Type]TypeStrategy{
		reflect.TypeOf(time.Duration(0)): DurationStrategy{},
		reflect.TypeOf(time.Time{}):      TimeStrategy{},
		reflect.TypeOf([]string{}):       StringSliceStrategy{},
		reflect.TypeOf([]bool{}):         BoolSliceStrategy{},
		reflect.TypeOf([]int{}):          IntSliceStrategy[int]{},
		reflect.TypeOf([]int8{}):         IntSliceStrategy[int8]{},
		reflect.TypeOf([]int16{}):        IntSliceStrategy[int16]{},
		reflect.TypeOf([]int32{}):        IntSliceStrategy[int32]{},
		reflect.TypeOf([]int64{}):        IntSliceStrategy[int64]{},
		reflect.TypeOf([]uint{}):         UintSliceStrategy[uint]{},
		reflect.TypeOf([]uint8{}):        UintSliceStrategy[uint8]{},
		reflect.TypeOf([]uint16{}):       UintSliceStrategy[uint16]{},
		reflect.TypeOf([]uint32{}):       UintSliceStrategy[uint32]{},
		reflect.TypeOf([]uint64{}):       UintSliceStrategy[uint64]{},
		reflect.TypeOf([]float64{}):      FloatSliceStrategy[float64]{},
		reflect.TypeOf([]float32{}):      FloatSliceStrategy[float32]{},
	}
}
