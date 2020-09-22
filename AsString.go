package config

import "strings"

func AsString() *StringChain {
	var chain StringChain
	empty := ""
	chain.empty = &empty
	return &chain
}

func (chain *StringChain) afterSetValue() {
	if chain.strval == nil && chain.value != nil {
		chain.strval = chain.value
	}
}

func (chain *StringChain) afterSetStringValue() {
	if chain.value == nil && chain.strval != nil {
		chain.value = chain.strval
	}
}

func (chain *StringChain) afterSetEmpty() {
	// nothing to do
}

func (chain *StringChain) trySetStringValue(value string) {

	// make sure it isn't empty
	value = strings.Trim(value, " ")
	if value == *chain.empty {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// set the value
	chain.value = &value

}

func (chain *StringChain) isEmpty(value string) bool {
	empty := *chain.empty
	return value == empty
}
