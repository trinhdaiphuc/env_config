package env_config

import (
	"fmt"
	"reflect"
	"time"
)

type TypeHandler interface {
	Handle(key string, field reflect.Value, nestedTagOpts TagOption) (Item, error)
}

type TypeHandlerFactory struct {
	handlers map[reflect.Type]TypeHandler
}

var (
	handlerFactory = NewTypeHandlerFactory()
)

func NewTypeHandlerFactory() *TypeHandlerFactory {
	return &TypeHandlerFactory{
		handlers: map[reflect.Type]TypeHandler{
			reflect.TypeOf(time.Time{}): TimeHandler{},
		},
	}
}

func (f *TypeHandlerFactory) GetHandler(t reflect.Type) TypeHandler {
	if handler, ok := f.handlers[t]; ok {
		return handler
	}

	// Default handler
	if t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct) {
		return StructHandler{}
	}

	return FieldHandler{}
}

type TimeHandler struct{}

func (h TimeHandler) Handle(key string, field reflect.Value, nestedTagOpt TagOption) (Item, error) {
	return FieldItem{
		raw:       field.Interface(),
		key:       key,
		value:     field,
		tagOption: nestedTagOpt,
	}, nil
}

type StructHandler struct{}

// Handle builds the nested section. Tag options other than the key are not
// applicable to a section and are intentionally ignored.
func (h StructHandler) Handle(key string, field reflect.Value, _ TagOption) (Item, error) {
	if field.Kind() != reflect.Ptr {
		return NewStruct(field.Addr().Interface(), key)
	}

	if field.Type().Elem().Kind() != reflect.Struct {
		return nil, fmt.Errorf("unsupported nested type %s for key %s", field.Type(), key)
	}

	// a section the caller already allocated is filled in place
	if !field.IsNil() {
		return NewStruct(field.Interface(), key)
	}

	alloc := reflect.New(field.Type().Elem())
	child, err := NewStruct(alloc.Interface(), key)
	if err != nil {
		return nil, err
	}
	child.target = field
	child.alloc = alloc

	return child, nil
}

type FieldHandler struct{}

func (h FieldHandler) Handle(key string, field reflect.Value, nestedTagOpt TagOption) (Item, error) {
	return FieldItem{
		raw:       field.Interface(),
		key:       key,
		value:     field,
		tagOption: nestedTagOpt,
	}, nil
}
