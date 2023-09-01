package config

import (
	"fmt"
	"strconv"
	"strings"
)

func AsFloat() *Float64Chain {
	var chain Float64Chain
	empty := 0.0
	chain.empty = &empty
	return &chain
}

func (chain *Float64Chain) afterSetValue() {
	// nothing to do
}

func (chain *Float64Chain) afterSetStringValue() {
	// nothing to do
}

func (chain *Float64Chain) afterSetEmpty() {
	// nothing to do
}

func (chain *Float64Chain) trySetStringValue(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// attempt to convert to int
	converted, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return
	}

	// set if not empty
	if converted != *chain.empty {
		chain.value = &converted
	}

}

func (chain *Float64Chain) isEmpty(value float64) bool {
	empty := *chain.empty
	return value == empty
}

func (chain *Float64Chain) Clamp(min float64, max float64) *Float64Chain {
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
