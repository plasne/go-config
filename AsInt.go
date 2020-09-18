package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type IntChain struct {
	name   *string
	strval *string
	value  *int
}

func (chain *IntChain) Name(name string) *IntChain {
	chain.name = &name
	return chain
}

func (chain *IntChain) setString(value string) {

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
	converted, err := strconv.Atoi(value)
	if err == nil {
		chain.value = &converted
	}

}

func (chain *IntChain) TrySetByEnv(key string) *IntChain {

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
func (chain *IntChain) Lookup(lookup map[string]int) *IntChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *int
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

func (chain *IntChain) TrySetByString(value string) *IntChain {
	chain.setString(value)
	return chain
}

func (chain *IntChain) SetTo(value int) *IntChain {
	chain.value = &value
	return chain
}

func (chain *IntChain) Clear() *IntChain {
	chain.value = nil
	return chain
}

func (chain *IntChain) TrySetTo(value int) *IntChain {
	if chain.value == nil {
		chain.value = &value
	}
	return chain
}

// Alias for TrySetTo().
func (chain *IntChain) DefaultTo(value int) *IntChain {
	return chain.TrySetTo(value)
}

func (chain *IntChain) Transform(f func(*IntChain)) *IntChain {
	f(chain)
	return chain
}

func (chain *IntChain) Print() *IntChain {
	fmt.Printf("  %s = %d\n", chain.Key(), chain.Value())
	return chain
}

func (chain *IntChain) PrintMasked() *IntChain {
	if chain.value != nil {
		fmt.Printf("  %s = (set)\n", chain.Key())
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
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

func (chain *IntChain) Require() *IntChain {
	if chain.value == nil {
		panic(fmt.Errorf("%s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *IntChain) IsSet() bool {
	return chain.value != nil
}

func (chain *IntChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *IntChain) Value() int {
	if chain.value == nil {
		return 0
	} else {
		return *chain.value
	}
}

func (chain *IntChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

func AsInt() *IntChain {
	var chain IntChain
	return &chain
}
