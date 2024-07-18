package env_config

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func Test_parseTag(t *testing.T) {
	type args struct {
		tag string
	}
	tests := []struct {
		name string
		args args
		want TagOption
	}{
		{
			name: "single key-value pair",
			args: args{
				tag: "default=value",
			},
			want: &DefaultOption{
				BaseTagOption: BaseTagOption{},
				DefaultValue:  "value",
			},
		},
		{
			name: "multiple key-value pairs",
			args: args{
				tag: "default=value1;delimiter=value2",
			},
			want: &DefaultOption{
				BaseTagOption: BaseTagOption{
					next: &DelimiterOption{
						BaseTagOption: BaseTagOption{},
						Delimiter:     "value2",
					},
				},
				DefaultValue: "value1",
			},
		},
		{
			name: "key with no value",
			args: args{
				tag: "default=",
			},
			want: &DefaultOption{
				BaseTagOption: BaseTagOption{},
				DefaultValue:  "",
			},
		},
		{
			name: "empty tag",
			args: args{
				tag: "",
			},
			want: nil,
		},
		{
			name: "multiple key-value pairs",
			args: args{
				tag: "default=value1,value2,value3;delimiter=,",
			},
			want: &DefaultOption{
				BaseTagOption: BaseTagOption{
					next: &DelimiterOption{
						BaseTagOption: BaseTagOption{},
						Delimiter:     ",",
					},
				},
				DefaultValue: "value1,value2,value3",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTag(tt.args.tag); !reflect.DeepEqual(got, tt.want) {
				fmt.Println(tt.args.tag)
				fmt.Println(strings.Split(tt.args.tag, Semicolon))
				fmt.Println(len(strings.SplitN(tt.args.tag, Equal, 2)))

				t.Errorf("parseTag() \ngot  = %#v \nwant = %#v", got, tt.want)
				t.Errorf("parseTag() \ngot  = %v \nwant = %v", toJSONStr(got), toJSONStr(tt.want))
			}
		})
	}
}
