package config

import (
)

type PropConfig struct {
    path string
    props map[string][string]
}

type ParseState int

const (
    KEY_START ParseState = iota
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
        make(map[string]string, 100),
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
    if val, ok = self.props[k]; ok {
        err = nil
    } else {
        err = fmt.Errorf("Don't contains key: \"%s\".", k)
    }
    return val, err
}

func (self *PropConfig) AddProp(k string, v interface{}) {
    self.props[k] = v.(string)
}

func (self *PropConfig) parseLine() {
}

func (self *PropConfig) load() error {
    var file, err = os.Open(self.path)
    if err != nil {
        return err
    }
    defer file.Close()

    var state = KEY_START
    var reader = bufio.NewReader(file)
    var curKey = ""
    var curVal = ""
    var curVar = ""
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
                // AddProp()
                // clear
            } else {
                // append to curKey
            }
        case KEY_ESCAPE:
            if c == '\n' {
                // addProp()
                // clear
                state = KEY_START
            } else {
                // append to curKey
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
                // append $
                state = KEY_START
            }
        case KEY_VAR:
            if c == '}' {
                // append to curKey with Get()
                // clear curvar
                state = KEY_START
            } else if c == '\n' {
                // append to curkey, contains ${
                // clear curvar
                state = KEY_START
            } else {
                // append to curval
            }
        case VAL_START:
            if c == '#' {
                state = VAL_COMMET
            } else if c == '\\' {
                state = VAL_ESCAPE
            } else if c == '$' {
                state = VAL_PREVAR
            } else if c == '\n' {
                // AddProp()
                // clear
                state = KEY_START
            } else {
                // append to curVal
            }
        case VAL_ESCAPE:
            if c == '\n' {
                state = VAL_START
            } else {
                // append to curVal
                state = VAL_START
            }
        case VAL_COMMENT:
            if c == '\n' {
                // addProp()
                // clear
                state = KEY_START
            } else {
            }
        case VAL_PREVAR:
            if c == '{' {
                state = VAL_VAR
            } else {
                // append $
                state = VAL_START
            }
        case VAL_VAR:
            if c == '}' {
                // append to curVal with Get()
                // clear curvar
                state = VAL_START
            } else if c == '\n' {
                // append to curvar, contains ${
                // clear curvar
                state = VAL_START
            } else {
                // append to curvar
            }
        case EOF:
            break
        }
    }

    return nil
}