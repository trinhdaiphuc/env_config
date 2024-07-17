package env_config

import (
	"reflect"
	"time"
)

type TypeHandler interface {
	Handle(key string, field reflect.Value, nestedTagOpts TagOption) Item
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

func (h TimeHandler) Handle(key string, field reflect.Value, nestedTagOpt TagOption) Item {
	return FieldItem{
		raw:       field.Interface(),
		key:       key,
		value:     field,
		tagOption: nestedTagOpt,
	}
}

type StructHandler struct{}

func (h StructHandler) Handle(key string, field reflect.Value, _ TagOption) Item {
	if field.Kind() == reflect.Ptr {
		field = field.Elem()
	}
	childStruct, err := NewStruct(field.Addr().Interface(), key)
	if err != nil {
		return nil
	}
	return childStruct
}

type FieldHandler struct{}

func (h FieldHandler) Handle(key string, field reflect.Value, nestedTagOpt TagOption) Item {
	return FieldItem{
		raw:       field.Interface(),
		key:       key,
		value:     field,
		tagOption: nestedTagOpt,
	}
}
