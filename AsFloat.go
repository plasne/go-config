package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FloatChain struct {
	name   *string
	strval *string
	value  *float64
}

func (chain *FloatChain) Name(name string) *FloatChain {
	chain.name = &name
	return chain
}

func (chain *FloatChain) setString(value string) {

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
	if err == nil {
		chain.value = &converted
	}

}

func (chain *FloatChain) TrySetByEnv(key string) *FloatChain {

	// set the name if not set
	if chain.name == nil {
		chain.name = &key
	}

	// ignore if already set
	if chain.value == nil {
		raw, ok := os.LookupEnv(key)
		if ok {
			chain.setString(raw)
		}
	}

	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *FloatChain) Lookup(lookup map[string]float64) *FloatChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *float64
		if v, ok := lookup[key]; ok {
			val = &v
		} else if v, ok := lookup[strings.ToLower(key)]; ok {
			val = &v
		}
		if val != nil {
			chain.value = val
		}
	}
	return chain
}

func (chain *FloatChain) TrySetByString(value string) *FloatChain {
	chain.setString(value)
	return chain
}

func (chain *FloatChain) SetTo(value float64) *FloatChain {
	chain.value = &value
	return chain
}

func (chain *FloatChain) Clear() *FloatChain {
	chain.value = nil
	return chain
}

func (chain *FloatChain) TrySetTo(value float64) *FloatChain {
	if chain.value == nil {
		chain.value = &value
	}
	return chain
}

// Alias for TrySetTo().
func (chain *FloatChain) DefaultTo(value float64) *FloatChain {
	return chain.TrySetTo(value)
}

func (chain *FloatChain) Transform(f func(*FloatChain)) *FloatChain {
	f(chain)
	return chain
}

func (chain *FloatChain) Print() *FloatChain {
	fmt.Printf("%s = %f\n", chain.Key(), chain.Value())
	return chain
}

func (chain *FloatChain) PrintMasked() *FloatChain {
	if chain.value != nil {
		fmt.Printf("%s = (set)\n", chain.Key())
	} else {
		fmt.Printf("%s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *FloatChain) Require() *FloatChain {
	if chain.value == nil {
		panic(fmt.Errorf("%s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *FloatChain) IsSet() bool {
	return chain.value != nil
}

func (chain *FloatChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *FloatChain) Value() float64 {
	if chain.value == nil {
		return 0
	} else {
		return *chain.value
	}
}

func (chain *FloatChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

func AsFloat() *FloatChain {
	var chain FloatChain
	return &chain
}
