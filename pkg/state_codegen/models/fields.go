package models

import (
	"strings"
	"unicode"
)

const (
	intType      = "int"
	stringType   = "string"
	durationType = "time.Duration"
)

type Field struct {
	SnakeName string `yaml:"name"`
	Type      string
}

func (f *Field) LowerName() string {
	name := f.Name()
	return string(unicode.ToLower(rune(name[0]))) + name[1:]
}

func (f *Field) Name() string {
	return toCamelInitCase(f.SnakeName)
}

func (f *Field) IsString() bool {
	return f.Type == stringType
}

func (f *Field) IsInt() bool {
	return f.Type == intType
}

func (f *Field) IsDuration() bool {
	return f.Type == durationType
}

type Params struct {
	Fields          []Field
	ClientName      string `yaml:"client_name"`
	ClientPath      string `yaml:"client_path"`
	PathToGenerated string `yaml:"path_to_generated"`
}

func (p Params) HasInt() bool {
	for _, el := range p.Fields {
		if el.Type == intType {
			return true
		}
	}
	return false
}

var uppercaseAcronym = map[string]string{
	"ID": "id",
}

// Converts a string to CamelCase
func toCamelInitCase(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	if a, ok := uppercaseAcronym[s]; ok {
		s = a
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := true
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
