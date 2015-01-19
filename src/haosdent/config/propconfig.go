package config

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "container/list"
)

type PropConfig struct {
    path string
    props map[string]*list.List
}

type PropParseState int

const (
    KEY_START PropParseState = iota
    KEY_ESCAPE
    KEY_COMMET
    KEY_PREVAR
    KEY_VAR

    VAL_START
    VAL_ESCAPE
    VAL_COMMENT
    VAL_PREVAR
    VAL_VAR

    EOF
)

func NewPropConfig(path string) *PropConfig {
    var instance = &PropConfig{
        path,
        make(map[string]*list.List, 100),
    }
    var err = instance.load()
    if err != nil {
        //panic(err)
        return nil
    } else {
        return instance
    }
}

func (self *PropConfig) Get(k string) (val string, err error) {
    var ok bool
    var l *list.List
    if l, ok = self.props[k]; ok {
        err = nil
    } else {
        err = fmt.Errorf("Don't contains key: \"%s\".", k)
    }
    val = strings.Join(l, ',')
    return val, err
}

func (self *PropConfig) AddProp(k string, v interface{}) {
    var ok bool
    var l *list.List
    if l, ok = self.props[k]; !ok {
        l = list.New()
    }
    l.PushBack(v.(string))
    self.props[k] = l;
}

func (self *PropConfig) load() error {
    var file, err = os.Open(self.path)
    if err != nil {
        return err
    }
    defer file.Close()

    var state = KEY_START
    var reader = bufio.NewReader(file)
    var curKey = make([]rune, 1024)
    var curKeyLen = 0
    var curVal = make([]rune, 1024)
    var curValLen = 0
    var curVar = make([]rune, 1024)
    var curVarLen = 0

    // TODO: when to switch state to EOF
    for state = EOF {
        var c = reader.ReadRune()
        switch state {
        case KEY_START:
            if c == '#' {
                state = KEY_COMMET
            } else if c == '\\' {
                state = KEY_ESCAPE
            } else if c == '$' {
                state = KEY_PREVAR
            } else if c == '\n' {
                // AddProp
                AddProp(string(curKey[:curKeyLen]), string(curVal[:curValLen]))
                // clear all
                curKeyLen = 0
                curValLen = 0
                curVarLen = 0
            } else {
                // append to curKey
                curKey[curKeyLen] = c
                curKeyLen++
            }
        case KEY_ESCAPE:
            if c == '\n' {
            } else {
                // append to curKey
                curKey[curKeyLen] = c
                curKeyLen++
                state = KEY_START
            }
        case KEY_COMMET:
            if c == '\n' {
                state = KEY_START
            } else {
            }
        case KEY_PREVAR:
            if c == '{' {
                state = KEY_VAR
            } else {
                // append '$' to curKey
                curKey[curKeyLen] = '$'
                curKeyLen++
                state = KEY_START
            }
        case KEY_VAR:
            if c == '}' {
                // append to curKey with Get()
                var tmpVal = Get(curVar)
                for _, ch := range tmpVal {
                    curKey[curKeyLen] = ch
                    curKeyLen++
                }
                // clear curVar
                curVarLen = 0
                state = KEY_START
            } else if c == '\n' {
                // append to curKey, contains ${
                curKey[curKeyLen] = '$'
                curKeyLen++
                curKey[curKeyLen] = '{'
                curKeyLen++
                for i := range curVarLen {
                    curKey[curKeyLen] = curVar[i]
                    curKeyLen++
                }
                // clear curVar
                curVarLen = 0
                state = KEY_START
            } else {
                // append to curVar
                curVar[curVarLen] = c
                curVarLen++
            }
        case VAL_START:
            if c == '#' {
                state = VAL_COMMET
            } else if c == '\\' {
                state = VAL_ESCAPE
            } else if c == '$' {
                state = VAL_PREVAR
            } else if c == ',' {
                // AddProp
                AddProp(string(curKey[:curKeyLen]), string(curVal[:curValLen]))
                // clear curVar and curVal
                curValLen = 0
                curVarLen = 0
            } else if c == '\n' {
                // AddProp
                AddProp(string(curKey[:curKeyLen]), string(curVal[:curValLen]))
                // clear all
                curKeyLen = 0
                curValLen = 0
                curVarLen = 0
                state = KEY_START
            } else {
                // append to curVal
                curVal[curValLen] = c
                curValLen++
            }
        case VAL_ESCAPE:
            if c == '\n' {
                state = VAL_START
            } else {
                // append to curVal
                curVal[curValLen] = c
                curValLen++
                state = VAL_START
            }
        case VAL_COMMENT:
            if c == '\n' {
                // AddProp
                AddProp(string(curKey[:curKeyLen]), string(curVal[:curValLen]))
                // clear all
                curKeyLen = 0
                curValLen = 0
                curVarLen = 0
                state = KEY_START
            } else {
            }
        case VAL_PREVAR:
            if c == '{' {
                state = VAL_VAR
            } else {
                // append '$' to curVal
                curVal[curValLen] = '$'
                curValLen++
                state = VAL_START
            }
        case VAL_VAR:
            if c == '}' {
                // append to curVal with Get()
                var tmpVal = Get(curVar)
                for _, ch := range tmpVal {
                    curVal[curValLen] = ch
                    curValLen++
                }
                // clear curVar
                curVarLen = 0
                state = VAL_START
            } else if c == '\n' {
                // append to curVal, contains ${
                curVal[curValLen] = '$'
                curValLen++
                curVal[curValLen] = '{'
                curValLen++
                for i := range curVarLen {
                    curVal[curValLen] = curVar[i]
                    curValLen++
                }
                // clear curVar
                curVarLen = 0
                state = VAL_START
            } else {
                // append to curVal
                curVal[curValLen] = c
                curValLen++
            }
        case EOF:
            break
        }
    }

    return nil
}

func (self *PropConfig) Save() error {
    var out, err = os.Create(self.path)
    if err != nil {
        return err
    }
    defer out.Close()

    for k, l := range self.props {
        for i, e := range l {
            l[i] = strings.Replace(e, ",", "\\,", -1)
        }
        var v = strings.Join(l, ", ")
        fmt.Fprintf(out, "%s = %s\n", k, v)
    }

    return nil
}