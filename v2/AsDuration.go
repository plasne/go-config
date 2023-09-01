package config

import (
	"strings"
	"time"
)

func AsDuration() *TimeDurationChain {
	var chain TimeDurationChain
	empty := time.Duration(0)
	chain.empty = &empty
	return &chain
}

func (chain *TimeDurationChain) afterSetValue() {
	// nothing to do
}

func (chain *TimeDurationChain) afterSetStringValue() {
	// nothing to do
}

func (chain *TimeDurationChain) afterSetEmpty() {
	// nothing to do
}

func (chain *TimeDurationChain) trySetStringValue(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// attempt to convert to duration
	converted, err := time.ParseDuration(value)
	if err != nil {
		return
	}

	// set if not empty
	if converted != *chain.empty {
		chain.value = &converted
	}

}

func (chain *TimeDurationChain) isEmpty(value time.Duration) bool {
	empty := *chain.empty
	return value.Seconds() == empty.Seconds()
}
