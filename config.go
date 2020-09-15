package config

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Convert string to int.
// Atoi(string, int) -OR- Atoi(string, func (string) int {})
func Atoi(original string, dflt interface{}) int {

	// attempt to convert to int
	converted, err := strconv.Atoi(original)
	if err == nil {
		return converted
	}

	// return default
	if dfltAsInt, ok := dflt.(int); ok {
		return dfltAsInt
	}
	if dfltAsFunc, ok := dflt.(func(string) int); ok {
		return dfltAsFunc(original)
	}
	return 0

}

func IsLower(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

// Convert string to int.
// Atoe(string, map[string]int, int) -OR- Atoe(string, map[string]int, func (string) int {})
func Atoe(original string, guide map[string]int, dflt interface{}) int {

	// ensure all map keys are lowercase
	for key := range guide {
		if !IsLower(key) {
			panic("the map in Atoe() must contain all lowercase keys")
		}
	}

	// see if the guide contains the key
	original = strings.ToLower(original)
	if val, ok := guide[original]; ok {
		return val
	}

	// return default
	if dfltAsInt, ok := dflt.(int); ok {
		return dfltAsInt
	}
	if dfltAsFunc, ok := dflt.(func(string) int); ok {
		return dfltAsFunc(original)
	}
	return 0

}

func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func Require(key string, value interface{}) {
	if valueAsString, ok := value.(string); ok {
		if len(valueAsString) > 0 {
			return
		}
		panic(fmt.Sprintf("%s is REQUIRED but missing.", key))
	}
	if valueAsInt, ok := value.(int); ok {
		if valueAsInt != 0 {
			return
		}
		panic(fmt.Sprintf("%s is REQUIRED but missing.", key))
	}
	panic(fmt.Sprintf("%s is REQUIRED but its value is an unrecognized type.", key))
}
