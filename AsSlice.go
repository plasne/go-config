package config

import (
	"strings"
)

type Slice []string

func AsSlice() *SliceChain {
	var chain SliceChain
	empty := Slice([]string{})
	chain.empty = &empty
	return &chain
}

func (chain *SliceChain) afterSetValue() {
	// nothing to do
}

func (chain *SliceChain) afterSetStringValue() {
	// nothing to do
}

func (chain *SliceChain) afterSetEmpty() {
	// nothing to do
}

func (chain *SliceChain) trySetStringValue(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// determine the delimiter
	delimiter := ","
	if chain.metadata != nil {
		if d, ok := chain.metadata["delimiter"].(string); ok {
			delimiter = d
		}
	}

	// split
	raw := strings.Split(value, delimiter)

	// normalize
	converted := make([]string, 0, len(raw))
	for _, v := range raw {
		t := strings.Trim(v, " ")
		if len(t) > 0 {
			converted = append(converted, t)
		}
	}

	// assign
	cast := Slice(converted)
	if !chain.isEmpty(cast) {
		chain.value = &cast
	}

}

func (chain *SliceChain) isEmpty(value Slice) bool {
	empty := *chain.empty
	return areSlicesEqual(value, empty)
}

func (chain *SliceChain) UseDelimiter(delimiter string) *SliceChain {
	if chain.metadata == nil {
		chain.metadata = make(map[string]interface{})
		chain.metadata["delimiter"] = delimiter
	}
	return chain
}
