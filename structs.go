package env_config

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (
	// DefaultTagName is the default tag name for struct fields which provides
	// a more granular to tweak certain structs. Lookup the necessary functions
	// for more info.
	DefaultTagName = "env" // struct field default tag name

)

var (
	_ Item = FieldItem{}
	_ Item = StructItem{}
)

type Item interface {
	TagOption() TagOption
	Value() reflect.Value
	Load() error
	Key() string
}

type FieldItem struct {
	raw       interface{}
	value     reflect.Value
	key       string
	tagOption TagOption
}

func (c FieldItem) Key() string {
	return c.key
}

func (c FieldItem) TagOption() TagOption {
	return c.tagOption
}

func (c FieldItem) Value() reflect.Value {
	return c.value
}

func (c FieldItem) Load() error {
	envValue := os.Getenv(c.key)

	// Ensure we have the correct kind of value to set
	value := c.value
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			// an optional field stays nil unless an env var or a default provides a value
			if !hasEnvValue(c) {
				return nil
			}
			if !value.CanSet() {
				return fmt.Errorf("cannot set value for key %s", c.key)
			}
			value.Set(reflect.New(value.Type().Elem()))
		}
		value = value.Elem()
	}

	if !value.CanSet() {
		return fmt.Errorf("cannot set value for key %s", c.key)
	}

	strategy, err := lookupStrategy(value.Type(), value.Kind())
	if err != nil {
		return fmt.Errorf("key %s: %w", c.key, err)
	}

	return strategy.SetValue(value, envValue, c.TagOption())
}

type StructItem struct {
	raw       interface{}
	prefix    string
	value     reflect.Value
	tagOption TagOption
	children  []Item

	// target is the nil pointer field this section fills; alloc is assigned to it
	// only when an env var or a default actually provides a value for a child.
	target reflect.Value
	alloc  reflect.Value
}

func (s StructItem) Load() error {
	if s.target.IsValid() && !hasEnvValue(s) {
		return nil
	}

	for _, child := range s.children {
		if err := child.Load(); err != nil {
			return err
		}
	}

	if !s.target.IsValid() {
		return nil
	}

	if !s.target.CanSet() {
		return fmt.Errorf("cannot set value for key %s", s.prefix)
	}
	s.target.Set(s.alloc)

	return nil
}

func (s StructItem) Key() string {
	return ""
}

func (s StructItem) TagOption() TagOption {
	return s.tagOption
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

		if !structField.IsExported() {
			return StructItem{}, fmt.Errorf("field %s: cannot load unexported field tagged with %q", structField.Name, DefaultTagName)
		}

		key, nestedTagOpts := parseTagAndKey(envTag)
		key = combineKeyPrefix(keyPrefix, key)

		fieldType := field.Type()
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		child, err := handlerFactory.GetHandler(fieldType).Handle(key, field, nestedTagOpts)
		if err != nil {
			return StructItem{}, fmt.Errorf("field %s: %w", structField.Name, err)
		}
		children = append(children, child)
	}

	return StructItem{
		prefix:   keyPrefix,
		raw:      s,
		value:    val,
		children: children,
	}, nil
}

// hasEnvValue reports whether any leaf under item resolves to a value, either
// from the environment or from a default tag option.
func hasEnvValue(item Item) bool {
	switch it := item.(type) {
	case FieldItem:
		if _, ok := os.LookupEnv(it.key); ok {
			return true
		}
		return hasDefaultValue(it.tagOption)
	case StructItem:
		for _, child := range it.children {
			if hasEnvValue(child) {
				return true
			}
		}
	}

	return false
}

func hasDefaultValue(option TagOption) bool {
	for ; option != nil; option = option.Next() {
		if defaultOpt, ok := option.(*DefaultOption); ok && defaultOpt != nil && defaultOpt.DefaultValue != "" {
			return true
		}
	}

	return false
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

func parseTagAndKey(str string) (key string, tag TagOption) {
	tagStr := strings.Split(str, Semicolon)
	if len(tagStr) <= 1 {
		key = str
		return
	}

	key = strings.TrimSpace(tagStr[0])
	tag = parseTag(strings.Join(tagStr[1:], Semicolon))
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
