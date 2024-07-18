package env_config

func LoadConfig(cfg interface{}) error {
	root, err := NewStruct(cfg, "")
	if err != nil {
		return err
	}
	return root.Load()
}
