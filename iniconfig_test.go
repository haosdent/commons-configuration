package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func SetUpIniConfigFile(path string) {
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

func TearDownIniConfigFile(path string) {
	var err = os.Remove(path)
	if err != nil {
		panic(err)
	}
}

func TestIniGet(t *testing.T) {
	var path = "config.ini"
	SetUpIniConfigFile(path)
	defer TearDownIniConfigFile(path)

	var c Configer = NewIniConfig(path)
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

func TestIniAddProp(t *testing.T) {
	var path = "config.ini"
	SetUpIniConfigFile(path)
	defer TearDownIniConfigFile(path)

	var c Configer = NewIniConfig(path)
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

func TestIniSave(t *testing.T) {
	var path = "config.ini"
	SetUpIniConfigFile(path)
	defer TearDownIniConfigFile(path)

	var c Configer = NewIniConfig(path)
	c.Save()

	var newC Configer = NewIniConfig(path)
	if !reflect.DeepEqual(c, newC) {
		t.Errorf("[after]Old struct and new struct is not equal.")
	}
}
