package env_config

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// DefaultTagName is the default tag name for struct fields which provides
	// a more granular to tweak certain structs. Lookup the necessary functions
	// for more info.
	DefaultTagName = "env" // struct field default tag name

	Underscore = "_"
)

var (
	_ Item = FieldItem{}
	_ Item = StructItem{}
)

type Item interface {
	TagOptions() []TagOption
	Value() reflect.Value
	Load() error
	Key() string
	print() string
}

type FieldItem struct {
	raw        interface{}
	value      reflect.Value
	key        string
	fieldName  string
	tagOptions []TagOption
}

func (c FieldItem) Key() string {
	return c.key
}

func (c FieldItem) print() string {
	str := "value:" + c.value.String() + ", key:" + c.key
	if c.tagOptions == nil {
		return str
	}
	str += " tagOptions:{"
	for _, tagOpt := range c.tagOptions {
		str += tagOpt.String()
	}

	return str + "}"
}

func (c FieldItem) TagOptions() []TagOption {
	return c.tagOptions
}

func (c FieldItem) Value() reflect.Value {
	return c.value
}

func (c FieldItem) Load() error {
	var defaultValue string
	for _, opt := range c.TagOptions() {
		if opt.key == DefaultTagKey {
			defaultValue = opt.value
		}
	}

	envValue := Env(c.key, defaultValue)

	// Ensure we have the correct kind of value to set
	value := c.value
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}

	if !value.CanSet() {
		return fmt.Errorf("cannot set value for key %s", c.key)
	}

	strategy, exists := complexTypeStrategies[value.Type()]
	if !exists {
		strategy, exists = buildInTypeStrategies[value.Kind()]
		if !exists {
			return nil
		}
	}

	return strategy.SetValue(value, envValue)
}

type StructItem struct {
	raw        interface{}
	prefix     string
	value      reflect.Value
	tagOptions []TagOption
	children   []Item
}

func (s StructItem) Load() error {
	for _, child := range s.children {
		if err := child.Load(); err != nil {
			return err
		}
	}
	return nil
}

func (s StructItem) Key() string {
	return ""
}

func (s StructItem) print() string {
	str := "value:" + s.value.String()
	if s.tagOptions != nil {
		str += " tagOptions:{"
		for _, tagOpt := range s.tagOptions {
			str += tagOpt.String()
		}
		str += "}"
	}

	str += "\nChildrens: "
	for _, child := range s.children {
		str += "\n" + child.print()
	}
	return str
}

func (s StructItem) TagOptions() []TagOption {
	return s.tagOptions
}

func (s StructItem) Value() reflect.Value {
	return s.value
}

func (s StructItem) Children() []Item {
	return s.children
}

func NewStruct(s interface{}, keyPrefix string) (StructItem, error) {
	val, err := pointerVal(s)
	if err != nil {
		return StructItem{}, err
	}

	typ := val.Type()

	var children []Item
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		structField := typ.Field(i)

		envTag := structField.Tag.Get(DefaultTagName)
		if envTag == "" {
			continue
		}

		envTag = strings.ReplaceAll(envTag, " ", "")
		key, nestedTagOpts := parseTagAndKey(envTag)
		key = combineKeyPrefix(keyPrefix, key)

		if field.Kind() == reflect.Ptr && field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}

		fieldType := field.Type()
		if field.Kind() == reflect.Ptr {
			fieldType = field.Elem().Type()
		}
		handler := handlerFactory.GetHandler(fieldType)
		children = append(children, handler.Handle(key, field, nestedTagOpts))
	}

	return StructItem{
		prefix:   keyPrefix,
		raw:      s,
		value:    val,
		children: children,
	}, nil
}

func pointerVal(s interface{}) (reflect.Value, error) {
	val := reflect.ValueOf(s)

	// if pointer get the underlying element≤
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return val, fmt.Errorf("expected struct, got %v", val.Kind())
	}
	return val, nil
}

func parseTagAndKey(str string) (key string, tags []TagOption) {
	tagStr := strings.Split(str, ",")
	if len(tagStr) <= 1 {
		key = str
		return
	}

	key = tagStr[0]
	tags = parseTag(strings.Join(tagStr[1:], ","))
	return
}

func combineKeyPrefix(prefix, key string) string {
	if prefix == "" {
		return key
	}
	if strings.HasSuffix(prefix, Underscore) {
		return prefix + key
	}

	return prefix + Underscore + key
}
