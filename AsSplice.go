package config

import (
	"fmt"
	"os"
	"strings"
)

type SpliceChain struct {
	name      *string
	strval    *string
	value     *[]string
	delimiter *string
}

func (chain *SpliceChain) Name(name string) *SpliceChain {
	chain.name = &name
	return chain
}

func (chain *SpliceChain) setString(value string) {

	// only proceed if there is a non-empty value
	value = strings.Trim(value, " ")
	if len(value) < 1 {
		return
	}

	// set if there is not already a strval
	if chain.strval == nil {
		chain.strval = &value
	}

	// determine the delimiter
	delimiter := ","
	if chain.delimiter != nil {
		delimiter = *chain.delimiter
	}

	// split
	raw := strings.Split(value, delimiter)

	// normalize
	converted := make([]string, 0, len(raw))
	for _, v := range raw {
		t := strings.Trim(v, " ")
		if len(t) > 0 {
			converted = append(converted, t)
		}
	}

	// assign
	if len(converted) > 0 {
		chain.value = &converted
	}

}

func (chain *SpliceChain) TrySetByEnv(key string) *SpliceChain {

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

func (chain *SpliceChain) TrySetByString(value string) *SpliceChain {
	chain.setString(value)
	return chain
}

func (chain *SpliceChain) SetTo(value []string) *SpliceChain {
	if value != nil {
		chain.value = &value
	} else {
		chain.value = nil
	}
	return chain
}

func (chain *SpliceChain) Clear() *SpliceChain {
	chain.value = nil
	return chain
}

func (chain *SpliceChain) TrySetTo(value []string) *SpliceChain {
	if chain.value == nil && value != nil {
		chain.value = &value
	}
	return chain
}

// Alias for TrySetTo().
func (chain *SpliceChain) DefaultTo(value []string) *SpliceChain {
	return chain.TrySetTo(value)
}

func (chain *SpliceChain) Transform(f func(*SpliceChain)) *SpliceChain {
	f(chain)
	return chain
}

func (chain *SpliceChain) Print() *SpliceChain {
	fmt.Printf("%s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *SpliceChain) PrintMasked() *SpliceChain {
	if chain.value != nil {
		fmt.Printf("%s = (set)\n", chain.Key())
	} else {
		fmt.Printf("%s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *SpliceChain) Require() *SpliceChain {
	if chain.value == nil {
		panic(fmt.Errorf("%s was REQUIRED but not provided.", chain.Key()))
	}
	return chain
}

func (chain *SpliceChain) UseDelimiter(delimiter string) *SpliceChain {
	chain.delimiter = &delimiter
	return chain
}

func (chain *SpliceChain) IsSet() bool {
	return chain.value != nil
}

func (chain *SpliceChain) Key() string {
	if chain.name != nil {
		return *chain.name
	} else {
		return "??????"
	}
}

func (chain *SpliceChain) Value() []string {
	if chain.value == nil {
		return []string{}
	} else {
		return *chain.value
	}
}

func AsSplice() *SpliceChain {
	var chain SpliceChain
	return &chain
}
