package config

import (
	"fmt"
	"strings"
)

func AsBool() *BoolChain {
	var chain BoolChain
	empty := false
	chain.empty = &empty
	return &chain
}

func (chain *BoolChain) afterSetValue() {
	// nothing to do
}

func (chain *BoolChain) afterSetStringValue() {
	// nothing to do
}

func (chain *BoolChain) afterSetEmpty() {
	panic(fmt.Errorf("SetEmpty() on a bool has no effect"))
}

func (chain *BoolChain) trySetStringValue(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// parse
	value = strings.ToLower(value)
	if value == "true" || value == "yes" || value == "y" || value == "1" {
		converted := true
		chain.value = &converted
	}
	if value == "false" || value == "no" || value == "n" || value == "0" {
		converted := false
		chain.value = &converted
	}

}

func (chain *BoolChain) isEmpty(value bool) bool {
	return false
}
