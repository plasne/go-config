package config

import (
	"fmt"
	"os"
	"strings"
)

type BoolChain struct {
	name   *string
	strval *string
	value  *bool
}

func (chain *BoolChain) Name(name string) *BoolChain {
	chain.name = &name
	return chain
}

func (chain *BoolChain) setString(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// attempt to convert to bool
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

func (chain *BoolChain) TrySetByEnv(key string) *BoolChain {

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
func (chain *BoolChain) Lookup(lookup map[string]bool) *BoolChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *bool
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

func (chain *BoolChain) TrySetByString(value string) *BoolChain {
	chain.setString(value)
	return chain
}

func (chain *BoolChain) SetTo(value bool) *BoolChain {
	chain.value = &value
	return chain
}

func (chain *BoolChain) Clear() *BoolChain {
	chain.value = nil
	return chain
}

func (chain *BoolChain) TrySetTo(value bool) *BoolChain {
	if chain.value == nil {
		chain.value = &value
	}
	return chain
}

// Alias for TrySetTo().
func (chain *BoolChain) DefaultTo(value bool) *BoolChain {
	return chain.TrySetTo(value)
}

func (chain *BoolChain) Transform(f func(*BoolChain)) *BoolChain {
	f(chain)
	return chain
}

func (chain *BoolChain) Print() *BoolChain {
	fmt.Printf("  %s = %t\n", chain.Key(), chain.Value())
	return chain
}

func (chain *BoolChain) PrintMasked() *BoolChain {
	if chain.value != nil {
		fmt.Printf("  %s = (set)\n", chain.Key())
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *BoolChain) Require() *BoolChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *BoolChain) IsSet() bool {
	return chain.value != nil
}

func (chain *BoolChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *BoolChain) Value() bool {
	if chain.value == nil {
		return false
	} else {
		return *chain.value
	}
}

func (chain *BoolChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

func AsBool() *BoolChain {
	var chain BoolChain
	return &chain
}
