package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/cheekybits/genny/generic"
)

//go:generate genny -in=$GOFILE -out=gen-$GOFILE gen "DataType=string,int,float64,bool,Slice,time.Duration"

type DataType generic.Type

type DataTypeChain struct {
	IChain
	key      *string
	strval   *string
	value    *DataType
	empty    *DataType
	metadata map[string]interface{}
}

func (chain *DataTypeChain) SetKey(key string) *DataTypeChain {
	chain.key = &key
	return chain
}

func (chain *DataTypeChain) SetStringValue(value string) *DataTypeChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *DataTypeChain) SetValue(value DataType) *DataTypeChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *DataTypeChain) SetEmpty(value DataType) *DataTypeChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *DataTypeChain) Clear() *DataTypeChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *DataTypeChain) TrySetValue(value DataType) *DataTypeChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *DataTypeChain) DefaultTo(value DataType) *DataTypeChain {
	return chain.TrySetValue(value)
}

func (chain *DataTypeChain) TrySetByEnv(key string) *DataTypeChain {

	// set the name if not set
	if chain.key == nil {
		chain.key = &key
	}

	// ignore if already set
	if chain.value == nil {
		raw, ok := os.LookupEnv(key)
		if ok {
			chain.trySetStringValue(raw)
		}
	}

	return chain
}

func (chain *DataTypeChain) TrySetByString(value string) *DataTypeChain {
	chain.trySetStringValue(value)
	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *DataTypeChain) Lookup(lookup map[string]DataType) *DataTypeChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *DataType
		if v, ok := lookup[key]; ok {
			val = &v
		} else if v, ok := lookup[strings.ToLower(key)]; ok {
			val = &v
		}
		if val != nil {
			chain.value = val
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *DataTypeChain) Transform(f func(*DataTypeChain)) *DataTypeChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *DataTypeChain) EnsureOneOf(options ...string) *DataTypeChain {

	// use the value or empty to evaluate
	strval := chain.StringValue()

	// look for a match
	found := false
	for i := 0; i < len(options); i++ {
		if options[i] == strval {
			found = true
		}
	}

	// if not found, clear strval and value
	if !found {
		chain.strval = nil
		chain.value = nil
	}

	return chain
}

func (chain *DataTypeChain) Resolve() *DataTypeChain {
	if chain.strval != nil {
		val, err := resolve(*chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *DataTypeChain) Print() *DataTypeChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *DataTypeChain) PrintMasked() *DataTypeChain {
	if chain.value != nil {
		fmt.Printf("  %s = (set)\n", chain.Key())
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *DataTypeChain) Require() *DataTypeChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *DataTypeChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *DataTypeChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *DataTypeChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *DataTypeChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *DataTypeChain) Value() DataType {
	if chain.value == nil {
		return *chain.empty
	} else {
		return *chain.value
	}
}

func (chain *DataTypeChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}
