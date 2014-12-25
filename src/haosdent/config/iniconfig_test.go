package config

import (
	"testing"
)

func TestGet(t *testing.T) {
	var path = "/tmp/config.ini"
	var c Configer = NewINIConfig(path)
	var key = "props.exist"

	var except = "true"
	// c.AddProp(key, except)
	var val, err = c.Get(key)
	if err != nil {
		t.Errorf("[after]Error is not empty.")
	}
	if val != except {
		t.Errorf("[after]Value \"%s\" is not equal to except \"%s\"", val, except)
	}

	key = "props.noexist"
	except = ""
	val, err = c.Get(key)
	if err == nil {
		t.Errorf("[after]Error is empty.")
	}
	if val != except {
		t.Errorf("[after]Value is not empty.")
	}
}

func TestAddProp(t *testing.T) {
	var path = "/tmp/config.ini"
	var c Configer = NewINIConfig(path)
	var key = "props.exist"

	var val, err = c.Get(key)
	if err == nil {
		t.Errorf("[before]Error is empty.")
	}
	if val != "" {
		t.Errorf("[before]Value is not empty.")
	}

	var except = "true"
	val = except
	c.AddProp(key, val)
	val, err = c.Get(key)
	if err != nil {
		t.Errorf("[after]Error is not empty.")
	}
	if val != except {
		t.Errorf("[after]Value \"%s\" is not equal to except \"%s\"", val, except)
	}
}
