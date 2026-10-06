package main

import (
	"strconv"
)

type ValueKind int
const (
	VALUE_KIND_NIL = iota
	VALUE_KIND_BOOL
	VALUE_KIND_NUMBER
	VALUE_KIND_STRING
)

type Value struct {
	kind ValueKind
	data any
}

func (value *Value) toString() string {
	switch value.kind {
		case VALUE_KIND_NIL:
			return "nil"
		case VALUE_KIND_BOOL:
			if value.data.(bool) {
				return "true"
			} else {
				return "false"
			}
		case VALUE_KIND_NUMBER: {
			return strconv.FormatFloat(value.data.(float64), 'f', -1, 64)
		}
		case VALUE_KIND_STRING: {
			return value.data.(string)
		}
	}

	panic("Unhandled value kind")
}
