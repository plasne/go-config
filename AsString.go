package config

import (
	"context"
	"fmt"
	"os"
	"strings"
)

type StringChain struct {
	name       *string
	value      *string
	allowEmpty bool
}

func (chain *StringChain) Name(name string) *StringChain {
	chain.name = &name
	return chain
}

func (chain *StringChain) setString(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if chain.allowEmpty || len(value) > 0 {
		chain.value = &value
	}

}

func (chain *StringChain) TrySetByEnv(key string) *StringChain {

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
func (chain *StringChain) Lookup(lookup map[string]string) *StringChain {
	if chain.value != nil {
		key := *chain.value
		var val *string
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

func (chain *StringChain) SetTo(value string) *StringChain {
	chain.value = &value
	return chain
}

func (chain *StringChain) Clear() *StringChain {
	chain.value = nil
	return chain
}

func (chain *StringChain) TrySetTo(value string) *StringChain {
	if chain.value == nil {
		chain.setString(value)
	}
	return chain
}

// Alias for TrySetTo().
func (chain *StringChain) DefaultTo(value string) *StringChain {
	return chain.TrySetTo(value)
}

func (chain *StringChain) Transform(f func(*StringChain)) *StringChain {
	f(chain)
	return chain
}

func (chain *StringChain) Resolve(ctx context.Context) *StringChain {
	if chain.value != nil {
		val, err := resolve(ctx, *chain.value)
		if err != nil {
			panic(err)
		}
		chain.value = &val
	}
	return chain
}

func (chain *StringChain) Print() *StringChain {
	fmt.Printf("  %s = \"%s\"\n", chain.Key(), chain.Value())
	return chain
}

func (chain *StringChain) PrintMasked() *StringChain {
	if chain.value != nil {
		fmt.Printf("  %s = (set)\n", chain.Key())
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *StringChain) Require() *StringChain {
	if chain.value == nil {
		panic(fmt.Errorf("%s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *StringChain) AllowEmpty() *StringChain {
	chain.allowEmpty = true
	return chain
}

func (chain *StringChain) IsSet() bool {
	return chain.value != nil
}

func (chain *StringChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *StringChain) Value() string {
	if chain.value == nil {
		return ""
	} else {
		return *chain.value
	}
}

func AsString() *StringChain {
	var chain StringChain
	return &chain
}
