package config

import (
	"fmt"
	"strconv"
	"strings"
)

func AsInt() *IntChain {
	var chain IntChain
	empty := 0
	chain.empty = &empty
	return &chain
}

func (chain *IntChain) afterSetValue() {
	// nothing to do
}

func (chain *IntChain) afterSetStringValue() {
	// nothing to do
}

func (chain *IntChain) afterSetEmpty() {
	// nothing to do
}

func (chain *IntChain) trySetStringValue(value string) {

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
	converted, err := strconv.Atoi(value)
	if err != nil {
		return
	}

	// set if not empty
	if converted != *chain.empty {
		chain.value = &converted
	}

}

func (chain *IntChain) isEmpty(value int) bool {
	empty := *chain.empty
	return value == empty
}

func (chain *IntChain) Clamp(min int, max int) *IntChain {
	if chain.value != nil {
		if max < min {
			panic(fmt.Errorf("max must be >= min"))
		}
		if *chain.value < min {
			chain.value = &min
		}
		if *chain.value > max {
			chain.value = &max
		}
	}
	return chain
}

func (chain *IntChain) PrintLookup(lookup map[string]int) *IntChain {
	val := chain.Value()
	for k, v := range lookup {
		if v == val {
			fmt.Printf("  %s = %s\n", chain.Key(), k)
			break
		}
	}
	return chain
}
