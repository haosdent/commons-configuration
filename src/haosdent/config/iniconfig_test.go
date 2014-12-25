package config

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	var path = "/tmp/config.ini"
	var c Configer = NewINIConfig(path)
	var key = "props.first"
	var val, _ = c.Get(key)
	fmt.Println(val)
}
