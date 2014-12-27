package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func SetUpConfigFile(path string) {
	var data = []byte(`
    global = a
    [props]
    exist = true
    `)
	var err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
		panic(err)
	}
}

func TearDownConfigFile(path string) {
	var err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestGet(t *testing.T) {
	var path = "config.ini"
	SetUpConfigFile(path)
	defer TearDownConfigFile(path)

	var c Configer = NewINIConfig(path)
	var key = "props.exist"

	var except = "true"
	// c.AddProp(key, except)
	var val, err = c.Get(key)
	if err != nil {
		t.Errorf("[after]Error is not empty: %s.", err)
	}
	if val != except {
		fmt.Println(len(val), len(except))
		for _, c := range val {
			fmt.Println(c)
		}
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
	var path = "config.ini"
	SetUpConfigFile(path)
	defer TearDownConfigFile(path)

	var c Configer = NewINIConfig(path)
	var key = "props.noexist"

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
		t.Errorf("[after]Error is not empty: %s.", err)
	}
	if val != except {
		t.Errorf("[after]Value \"%s\" is not equal to except \"%s\"", val, except)
	}
}

func TestSave(t *testing.T) {
	var path = "config.ini"
	SetUpConfigFile(path)
	defer TearDownConfigFile(path)

	var c Configer = NewINIConfig(path)
	c.Save()

	var newC Configer = NewINIConfig(path)
	if !reflect.DeepEqual(c, newC) {
		t.Errorf("[after]Old struct and new struct is not equal.")
	}
}
