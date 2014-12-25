package config

type INIConfig struct {
	path  string
	props map[string]string
}

func NewINIConfig(path string) *INIConfig {
	var props = make(map[string]string, 100)
	props["props.first"] = "true"
	return &INIConfig{
		path,
		props,
	}
}

func (self *INIConfig) Get(k string) (val string, err error) {
	val = self.props[k]
	err = nil
	return val, err
}

func (self *INIConfig) AddProp(k string, v interface{}) {
}
