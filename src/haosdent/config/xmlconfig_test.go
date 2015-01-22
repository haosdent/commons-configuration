package config

import (
    "fmt"
    "io/ioutil"
    "os"
    "reflect"
    "testing"
)

func SetUpXmlConfigFile(path string) {
    var data = []byte(`
    <global>
      a
    </global>
    <props>
      <exist>true</exist>
    </props>
    `)
    var err = ioutil.WriteFile(path, data, 0644)
    if err != nil {
        panic(err)
    }
}

func TearDownXmlConfigFile(path string) {
    var err = os.Remove(path)
    if err != nil {
        panic(err)
    }
}