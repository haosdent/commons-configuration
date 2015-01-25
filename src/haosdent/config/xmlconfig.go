package config

import (
    "encoding/xml"
    "fmt"
    "os"
    "strings"
)

type XmlConfig struct {
    path  string
    props map[string][]string
}

func NewXmlConfig(path string) *XmlConfig {
    var instance = &XmlConfig{
        path,
        make(map[string][]string, 100),
    }
    var err = instance.load()
    if err != nil {
        //panic(err)
        return nil
    } else {
        return instance
    }
}

func (self *XmlConfig) Get(k string) (val string, err error) {
    var ok bool
    var l []string
    if l, ok = self.props[k]; ok {
        err = nil
        val = strings.Join(l, ",")
    } else {
        err = fmt.Errorf("Don't contains key: \"%s\".", k)
    }
    return val, err
}

func (self *XmlConfig) AddProp(k string, v interface{}) {
    var ok bool
    var l []string
    if l, ok = self.props[k]; !ok {
        l = make([]string, 0, 3)
    }
    k = strings.TrimSpace(k)
    if (len(k)) == 0 {
        return
    }

    var vStr = strings.TrimSpace(v.(string))
    l = append(l, vStr);
    self.props[k] = l;
}

func (self *XmlConfig) load() error {
    var file, err = os.Open(self.path)
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := xml.NewDecoder(file)

    var elems []string
    for {
        token, _ := decoder.Token()
        if token == nil {
            break
        }
        switch elem := token.(type) {
        case xml.StartElement:
            elems = append(elems, elem.Name.Local)
        case xml.EndElement:
            elems = elems[:len(elems) - 1]
        case xml.CharData:
            content := string([]byte(elem))
            content = strings.TrimSpace(content)
            if len(content) > 0 {
                self.AddProp(strings.Join(elems, "."), content)
            }
        }
    }

    return nil
}

func (self *XmlConfig) Save() error {
    return nil
}