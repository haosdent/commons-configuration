package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type INIConfig struct {
	path  string
	props map[string]string
}

type ParseState int

const (
	START ParseState = iota
	NORMAL
	QUOTED
	ESCAPE
	STOP
)

func NewINIConfig(path string) *INIConfig {
	var instance = &INIConfig{
		path,
		make(map[string]string, 100),
	}
	instance.load()
	return instance
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

func (self *INIConfig) isCommentLine(line string) bool {
	if line == "" {
		return false
	}

	return line[0] == '#' || line[0] == ';'
}

func (self *INIConfig) isSectionLine(line string) bool {
	if line == "" {
		return false
	}

	return line[0] == '[' && line[len(line)-1] == ']'
}

func (self *INIConfig) isSpace(c rune) bool {
	if c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r' {
		return true
	} else {
		return false
	}
}

func (self *INIConfig) parseVal(val string) string {
	if len(val) == 0 {
		return val
	}

	var state = START
	var bs = make([]rune, len(val))
	var bl = 0
	var quote rune

	for _, c := range val {
		switch state {
		case START:
			if c == '"' || c == '\'' {
				quote = c
				state = QUOTED
			} else if c == '#' || c == ';' {
				state = STOP
			} else if !self.isSpace(c) {
				bs[bl] = c
				bl++
				state = NORMAL
			}
		case NORMAL:
			if c == '#' || c == ';' || self.isSpace(c) {
				state = STOP
			} else {
				bs[bl] = c
				bl++
			}
		case QUOTED:
			if c == '\\' {
				state = ESCAPE
			} else if c == quote {
				state = STOP
			} else {
				bs[bl] += c
				bl++
			}
		case ESCAPE:
			if c != quote {
				bs[bl] += '\\'
				bl++
			}
			bs[bl] += c
			bl++
			state = QUOTED
		case STOP:
			break
		}
	}

	val = string(bs[:bl])

	return val
}

func (self *INIConfig) load() {
	var file, _ = os.Open(self.path)
	defer file.Close()

	var scanner = bufio.NewScanner(file)
	var section = ""
	for scanner.Scan() {
		var line = scanner.Text()
		line = strings.TrimSpace(line)
		if self.isCommentLine(line) || len(line) == 0 {
			continue
		} else if self.isSectionLine(line) {
			section = line[1:len(line)-1] + "."
			continue
		}

		var index = strings.Index(line, "=")
		if index < 0 {
			index = strings.Index(line, ":")
		}

		if index < 0 {
			var key = section + line
			self.props[key] = ""
		} else {
			var key = line[:index]
			key = section + strings.TrimSpace(key)
			var val = line[index+1:]
			val = self.parseVal(val)
			self.props[key] = val
			//fmt.Println(key, val)
		}
	}
}
