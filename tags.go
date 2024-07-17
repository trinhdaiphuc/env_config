package env_config

import (
	"sort"
	"strings"
)

// TagOption interface
type TagOption interface {
	Apply(value string) (interface{}, error)
	SetValue(value string)
	SetNext(option TagOption)
}

type TagOptionPriority interface {
	Priority() int
}

type TagOptionBuilder interface {
	Build() TagOption
}

const (
	DefaultTagKey = "default"
	Delimiter     = "delimiter"
)

const (
	Underscore = "_"
	Semicolon  = ";"
	Comma      = ","
	Colon      = ":"
	Equal      = "="
)

var (
	tagOptionBuilders = map[string]TagOptionBuilder{
		DefaultTagKey: &DefaultOptionBuilder{},
		Delimiter:     &DelimiterOptionBuilder{},
	}
)

type DefaultOptionBuilder struct{}

func (d *DefaultOptionBuilder) Build() TagOption {
	return &DefaultOption{
		BaseTagOption: BaseTagOption{},
		DefaultValue:  "",
	}
}

type DelimiterOptionBuilder struct{}

func (d *DelimiterOptionBuilder) Build() TagOption {
	return &DelimiterOption{
		BaseTagOption: BaseTagOption{},
		Delimiter:     "",
	}
}

// BaseTagOption to hold the next TagOption in the chain
type BaseTagOption struct {
	next TagOption
}

func (b *BaseTagOption) SetNext(option TagOption) {
	if b == nil {
		return
	}
	b.next = option
}

func (b *BaseTagOption) Apply(value string) (interface{}, error) {
	if b == nil || b.next == nil {
		return value, nil
	}

	return b.next.Apply(value)
}

// DefaultOption implementation
type DefaultOption struct {
	BaseTagOption
	DefaultValue string
}

func (d *DefaultOption) SetValue(value string) {
	d.DefaultValue = value
}

func (d *DefaultOption) Apply(value string) (interface{}, error) {
	if d == nil {
		return value, nil
	}
	if value == "" {
		value = d.DefaultValue
	}
	return d.BaseTagOption.Apply(value)
}

func (d *DefaultOption) Priority() int {
	return 0
}

// DelimiterOption implementation
type DelimiterOption struct {
	BaseTagOption
	Delimiter string
}

func (d *DelimiterOption) SetValue(value string) {
	d.Delimiter = value
}

func (d *DelimiterOption) Apply(value string) (interface{}, error) {
	if d == nil {
		return value, nil
	}

	if d.Delimiter == "" {
		d.Delimiter = ","
	}

	return strings.Split(value, d.Delimiter), nil
}

func (d *DelimiterOption) Priority() int {
	return 1
}

func parseTag(tag string) TagOption {
	parts := strings.Split(tag, Semicolon)
	var (
		head, tail  TagOption
		tempOptions []TagOptionPriority
	)

	for _, tag := range parts {
		parts := strings.SplitN(tag, "=", 2)
		if len(parts) == 2 {
			builder, ok := tagOptionBuilders[parts[0]]
			if !ok {
				continue
			}
			option := builder.Build()
			option.SetValue(parts[1])
			tempOptions = append(tempOptions, option.(TagOptionPriority))
		}

	}

	sort.Slice(tempOptions, func(i, j int) bool {
		return tempOptions[i].Priority() < tempOptions[j].Priority()
	})

	for _, tempOption := range tempOptions {
		opt := tempOption.(TagOption)
		if head == nil {
			head = opt
			tail = opt
		} else {
			tail.SetNext(opt)
			tail = opt
		}
	}

	if head == nil {
		builder := DefaultOptionBuilder{}
		head = builder.Build()
	}

	return head
}
