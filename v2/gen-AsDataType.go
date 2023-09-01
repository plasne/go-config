// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

package config

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
)

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "String=string,int,float64,bool,Slice,time.Duration"

type StringChain struct {
	IChain
	key      *string
	strval   *string
	value    *string
	empty    *string
	metadata map[string]interface{}
}

func (chain *StringChain) SetKey(key string) *StringChain {
	chain.key = &key
	return chain
}

func (chain *StringChain) SetStringValue(value string) *StringChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *StringChain) SetValue(value string) *StringChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *StringChain) SetEmpty(value string) *StringChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *StringChain) Clear() *StringChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *StringChain) TrySetValue(value string) *StringChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *StringChain) DefaultTo(value string) *StringChain {
	return chain.TrySetValue(value)
}

func (chain *StringChain) TrySetByEnv(key string) *StringChain {

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

func (chain *StringChain) TrySetByString(value string) *StringChain {
	chain.trySetStringValue(value)
	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *StringChain) Lookup(lookup map[string]string) *StringChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *string
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

func (chain *StringChain) Transform(f func(*StringChain)) *StringChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *StringChain) EnsureOneOf(options ...string) *StringChain {

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

func (chain *StringChain) Resolve(ctx context.Context) *StringChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *StringChain) Print() *StringChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *StringChain) PrintMasked() *StringChain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *StringChain) Require() *StringChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *StringChain) RequireIf(clause bool) *StringChain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *StringChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *StringChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *StringChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *StringChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *StringChain) Value() string {
	if chain.value == nil {
		return *chain.empty
	} else {
		return *chain.value
	}
}

func (chain *StringChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "Int=string,int,float64,bool,Slice,time.Duration"

type IntChain struct {
	IChain
	key      *string
	strval   *string
	value    *int
	empty    *int
	metadata map[string]interface{}
}

func (chain *IntChain) SetKey(key string) *IntChain {
	chain.key = &key
	return chain
}

func (chain *IntChain) SetStringValue(value string) *IntChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *IntChain) SetValue(value int) *IntChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *IntChain) SetEmpty(value int) *IntChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *IntChain) Clear() *IntChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *IntChain) TrySetValue(value int) *IntChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *IntChain) DefaultTo(value int) *IntChain {
	return chain.TrySetValue(value)
}

func (chain *IntChain) TrySetByEnv(key string) *IntChain {

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

func (chain *IntChain) TrySetByString(value string) *IntChain {
	chain.trySetStringValue(value)
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
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *IntChain) Transform(f func(*IntChain)) *IntChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *IntChain) EnsureOneOf(options ...string) *IntChain {

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

func (chain *IntChain) Resolve(ctx context.Context) *IntChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *IntChain) Print() *IntChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *IntChain) PrintMasked() *IntChain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *IntChain) Require() *IntChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *IntChain) RequireIf(clause bool) *IntChain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *IntChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *IntChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *IntChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *IntChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *IntChain) Value() int {
	if chain.value == nil {
		return *chain.empty
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

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "Float64=string,int,float64,bool,Slice,time.Duration"

type Float64Chain struct {
	IChain
	key      *string
	strval   *string
	value    *float64
	empty    *float64
	metadata map[string]interface{}
}

func (chain *Float64Chain) SetKey(key string) *Float64Chain {
	chain.key = &key
	return chain
}

func (chain *Float64Chain) SetStringValue(value string) *Float64Chain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *Float64Chain) SetValue(value float64) *Float64Chain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *Float64Chain) SetEmpty(value float64) *Float64Chain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *Float64Chain) Clear() *Float64Chain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *Float64Chain) TrySetValue(value float64) *Float64Chain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *Float64Chain) DefaultTo(value float64) *Float64Chain {
	return chain.TrySetValue(value)
}

func (chain *Float64Chain) TrySetByEnv(key string) *Float64Chain {

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

func (chain *Float64Chain) TrySetByString(value string) *Float64Chain {
	chain.trySetStringValue(value)
	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *Float64Chain) Lookup(lookup map[string]float64) *Float64Chain {
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
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *Float64Chain) Transform(f func(*Float64Chain)) *Float64Chain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *Float64Chain) EnsureOneOf(options ...string) *Float64Chain {

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

func (chain *Float64Chain) Resolve(ctx context.Context) *Float64Chain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *Float64Chain) Print() *Float64Chain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *Float64Chain) PrintMasked() *Float64Chain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *Float64Chain) Require() *Float64Chain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *Float64Chain) RequireIf(clause bool) *Float64Chain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *Float64Chain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *Float64Chain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *Float64Chain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *Float64Chain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *Float64Chain) Value() float64 {
	if chain.value == nil {
		return *chain.empty
	} else {
		return *chain.value
	}
}

func (chain *Float64Chain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "Bool=string,int,float64,bool,Slice,time.Duration"

type BoolChain struct {
	IChain
	key      *string
	strval   *string
	value    *bool
	empty    *bool
	metadata map[string]interface{}
}

func (chain *BoolChain) SetKey(key string) *BoolChain {
	chain.key = &key
	return chain
}

func (chain *BoolChain) SetStringValue(value string) *BoolChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *BoolChain) SetValue(value bool) *BoolChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *BoolChain) SetEmpty(value bool) *BoolChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *BoolChain) Clear() *BoolChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *BoolChain) TrySetValue(value bool) *BoolChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *BoolChain) DefaultTo(value bool) *BoolChain {
	return chain.TrySetValue(value)
}

func (chain *BoolChain) TrySetByEnv(key string) *BoolChain {

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

func (chain *BoolChain) TrySetByString(value string) *BoolChain {
	chain.trySetStringValue(value)
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
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *BoolChain) Transform(f func(*BoolChain)) *BoolChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *BoolChain) EnsureOneOf(options ...string) *BoolChain {

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

func (chain *BoolChain) Resolve(ctx context.Context) *BoolChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *BoolChain) Print() *BoolChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *BoolChain) PrintMasked() *BoolChain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *BoolChain) Require() *BoolChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *BoolChain) RequireIf(clause bool) *BoolChain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *BoolChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *BoolChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *BoolChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *BoolChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *BoolChain) Value() bool {
	if chain.value == nil {
		return *chain.empty
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

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "Slice=string,int,float64,bool,Slice,time.Duration"

type SliceChain struct {
	IChain
	key      *string
	strval   *string
	value    *Slice
	empty    *Slice
	metadata map[string]interface{}
}

func (chain *SliceChain) SetKey(key string) *SliceChain {
	chain.key = &key
	return chain
}

func (chain *SliceChain) SetStringValue(value string) *SliceChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *SliceChain) SetValue(value Slice) *SliceChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *SliceChain) SetEmpty(value Slice) *SliceChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *SliceChain) Clear() *SliceChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *SliceChain) TrySetValue(value Slice) *SliceChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *SliceChain) DefaultTo(value Slice) *SliceChain {
	return chain.TrySetValue(value)
}

func (chain *SliceChain) TrySetByEnv(key string) *SliceChain {

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

func (chain *SliceChain) TrySetByString(value string) *SliceChain {
	chain.trySetStringValue(value)
	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *SliceChain) Lookup(lookup map[string]Slice) *SliceChain {
	if chain.strval != nil {
		key := *chain.strval
		var val *Slice
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

func (chain *SliceChain) Transform(f func(*SliceChain)) *SliceChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *SliceChain) EnsureOneOf(options ...string) *SliceChain {

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

func (chain *SliceChain) Resolve(ctx context.Context) *SliceChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *SliceChain) Print() *SliceChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *SliceChain) PrintMasked() *SliceChain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *SliceChain) Require() *SliceChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *SliceChain) RequireIf(clause bool) *SliceChain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *SliceChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *SliceChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *SliceChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *SliceChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *SliceChain) Value() Slice {
	if chain.value == nil {
		return *chain.empty
	} else {
		return *chain.value
	}
}

func (chain *SliceChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}

//go:generate go run github.com/cheekybits/genny -in=$GOFILE -out=gen-$GOFILE gen "TimeDuration=string,int,float64,bool,Slice,time.Duration"

type TimeDurationChain struct {
	IChain
	key      *string
	strval   *string
	value    *time.Duration
	empty    *time.Duration
	metadata map[string]interface{}
}

func (chain *TimeDurationChain) SetKey(key string) *TimeDurationChain {
	chain.key = &key
	return chain
}

func (chain *TimeDurationChain) SetStringValue(value string) *TimeDurationChain {
	chain.strval = &value
	chain.afterSetStringValue()
	return chain
}

func (chain *TimeDurationChain) SetValue(value time.Duration) *TimeDurationChain {
	chain.value = &value
	chain.afterSetValue()
	return chain
}

func (chain *TimeDurationChain) SetEmpty(value time.Duration) *TimeDurationChain {
	chain.empty = &value
	chain.afterSetEmpty()
	return chain
}

func (chain *TimeDurationChain) Clear() *TimeDurationChain {
	chain.value = nil
	chain.afterSetValue()
	return chain
}

func (chain *TimeDurationChain) TrySetValue(value time.Duration) *TimeDurationChain {
	if chain.value == nil {
		if !chain.isEmpty(value) {
			chain.value = &value
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *TimeDurationChain) DefaultTo(value time.Duration) *TimeDurationChain {
	return chain.TrySetValue(value)
}

func (chain *TimeDurationChain) TrySetByEnv(key string) *TimeDurationChain {

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

func (chain *TimeDurationChain) TrySetByString(value string) *TimeDurationChain {
	chain.trySetStringValue(value)
	return chain
}

// Lookup() will replace the current value even if set if a suitable key/value pair was found.
func (chain *TimeDurationChain) Lookup(lookup map[string]time.Duration) *TimeDurationChain {
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
			chain.afterSetValue()
		}
	}
	return chain
}

func (chain *TimeDurationChain) Transform(f func(*TimeDurationChain)) *TimeDurationChain {
	f(chain)
	return chain
}

// EnsureOneOf() clears strval and value if strval is not one of the selected options.
func (chain *TimeDurationChain) EnsureOneOf(options ...string) *TimeDurationChain {

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

func (chain *TimeDurationChain) Resolve(ctx context.Context) *TimeDurationChain {
	if chain.strval != nil {
		val, err := resolve(ctx, *chain.strval)
		if err != nil {
			panic(err)
		}
		chain.trySetStringValue(val)
	}
	return chain
}

func (chain *TimeDurationChain) Print() *TimeDurationChain {
	fmt.Printf("  %s = %v\n", chain.Key(), chain.Value())
	return chain
}

func (chain *TimeDurationChain) PrintMasked() *TimeDurationChain {
	if chain.value != nil {
		if v := strings.ToLower(*chain.strval); strings.HasPrefix(v, "https://") && strings.Contains(v, ".vault.azure.net") {
			return chain.Print()
		} else {
			fmt.Printf("  %s = (set)\n", chain.Key())
		}
	} else {
		fmt.Printf("  %s = (not-set)\n", chain.Key())
	}
	return chain
}

func (chain *TimeDurationChain) Require() *TimeDurationChain {
	if chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *TimeDurationChain) RequireIf(clause bool) *TimeDurationChain {
	if clause && chain.value == nil {
		panic(fmt.Errorf("  %s was REQUIRED but not provided", chain.Key()))
	}
	return chain
}

func (chain *TimeDurationChain) IsKeySet() bool {
	return chain.key != nil
}

func (chain *TimeDurationChain) IsStringValueSet() bool {
	return chain.strval != nil
}

func (chain *TimeDurationChain) IsValueSet() bool {
	return chain.value != nil
}

func (chain *TimeDurationChain) Key() string {
	if chain.key != nil {
		return *chain.key
	} else {
		return "??????"
	}
}

func (chain *TimeDurationChain) Value() time.Duration {
	if chain.value == nil {
		return *chain.empty
	} else {
		return *chain.value
	}
}

func (chain *TimeDurationChain) StringValue() string {
	if chain.strval == nil {
		return ""
	} else {
		return *chain.strval
	}
}
