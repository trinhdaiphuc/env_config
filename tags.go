package env_config

import "strings"

const (
	DefaultTagKey = "default"
)

type TagOption struct {
	key   string
	value string
}

func (t TagOption) String() string {
	return t.key + ":" + t.value
}

func parseTag(tag string) []TagOption {
	parts := strings.Split(tag, ",")
	opts := make([]TagOption, len(parts))
	for i, part := range parts {
		kv := strings.SplitN(part, "=", 2)
		if len(kv) == 2 {
			opts[i] = TagOption{key: kv[0], value: kv[1]}
		} else {
			opts[i] = TagOption{key: kv[0], value: ""}
		}
	}
	return opts
}
