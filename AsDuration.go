package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

type DurationChain struct {
	name   *string
	strval *string
	value  *time.Duration
}

func (chain *DurationChain) Name(name string) *DurationChain {
	chain.name = &name
	return chain
}

func (chain *DurationChain) setString(value string) {

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
	if err == nil {
		chain.value = &converted
	}

}

func (chain *DurationChain) TrySetByEnv(key string) *DurationChain {

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
func (chain *DurationChain) Lookup(lookup map[string]time.Duration) *DurationChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *time.Duration
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

func (chain *DurationChain) TrySetByString(value string) *DurationChain {
	chain.setString(value)
	return chain
}

func (chain *DurationChain) SetTo(value time.Duration) *DurationChain {
	chain.value = &value
	return chain
}

func (chain *DurationChain) Clear() *DurationChain {
	chain.value = nil
	return chain
}

func (chain *DurationChain) TrySetTo(value time.Duration) *DurationChain {
	if chain.value == nil {
		chain.value = &value
	}
	return chain
}

// Alias for TrySetTo().
func (chain *DurationChain) DefaultTo(value time.Duration) *DurationChain {
	return chain.TrySetTo(value)
}

func (chain *DurationChain) Transform(f func(*DurationChain)) *DurationChain {
	f(chain)
	return chain
}

func (chain *DurationChain) Resolve(ctx context.Context) *DurationChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.setString(val)
	}
	return chain
}

func (chain *DurationChain) Print() *DurationChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *DurationChain) PrintMasked() *DurationChain {
	if chain.value != nil {
		fmt.Printf("  %s = (set)\n", chain.Key())
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *DurationChain) Require() *DurationChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *DurationChain) IsSet() bool {
	return chain.value != nil
}

func (chain *DurationChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *DurationChain) Value() time.Duration {
	if chain.value == nil {
		return time.Duration(0)
	} else {
		return *chain.value
	}
}

func (chain *DurationChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

func AsDuration() *DurationChain {
	var chain DurationChain
	return &chain
}
