package config

import (
	"fmt"
)

type INIConfig struct {
	path  string
	props map[string]string
}

func NewINIConfig(path string) *INIConfig {
	var props = make(map[string]string, 100)
	return &INIConfig{
		path,
		props,
	}
}

func (self *INIConfig) Get(k string) (val string, err error) {
	var ok bool
	if val, ok = self.props[k]; ok {
		err = nil
	} else {
		err = fmt.Errorf("Don't contains key: \"%s\".", k)
	}
	return val, err
}

func (self *INIConfig) AddProp(k string, v interface{}) {
	self.props[k] = v.(string)
}
