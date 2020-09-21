package config

import (
	"os"
	"testing"
	"time"
)

func TestAsString(t *testing.T) {

	t.Run("AsString().TrySetValue(cat)", func(t *testing.T) {
		e := "cat"
		a := AsString().TrySetValue("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetValue().DefaultTo(cat)", func(t *testing.T) {
		e := "cat"
		a := AsString().TrySetValue("").DefaultTo("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().SetEmpty(null).TrySetValue().DefaultTo(cat)", func(t *testing.T) {
		e := ""
		a := AsString().SetEmpty("null").TrySetValue("").DefaultTo("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetByEnv(dog)", func(t *testing.T) {
		e := "dog"
		os.Setenv("TEST_VALUE", "dog")
		a := AsString().TrySetByEnv("TEST_VALUE").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetByEnv().DefaultTo(dog)", func(t *testing.T) {
		e := "dog"
		os.Setenv("TEST_VALUE", "")
		a := AsString().TrySetByEnv("TEST_VALUE").DefaultTo("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().SetEmpty(null).TrySetByEnv().DefaultTo(dog)", func(t *testing.T) {
		e := ""
		os.Setenv("TEST_VALUE", "")
		a := AsString().SetEmpty("null").TrySetByEnv("TEST_VALUE").DefaultTo("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetByEnv(cat).SetValue(dog)", func(t *testing.T) {
		e := "dog"
		os.Setenv("TEST_VALUE", "cat")
		a := AsString().TrySetByEnv("TEST_VALUE").SetValue("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().SetValue(cat).SetValue(dog)", func(t *testing.T) {
		e := "dog"
		a := AsString().SetValue("cat").SetValue("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString()_null_is_empty", func(t *testing.T) {
		e := ""
		a := AsString().Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString()_whitespace_is_allowed", func(t *testing.T) {
		e := "   "
		a := AsString().TrySetValue("   ").DefaultTo("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString()_empty_is_not_set_on_try", func(t *testing.T) {
		a := AsString().TrySetValue("")
		if a.IsValueSet() {
			t.Errorf("AsString() Failed: expected IsValueSet() to be false")
		}
	})

	t.Run("AsString()_empty_is_set_on_explicit", func(t *testing.T) {
		a := AsString().SetValue("")
		if !a.IsValueSet() {
			t.Errorf("AsString() Failed: expected IsValueSet() to be true")
		}
	})

	t.Run("AsString()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsString() Failed: expected panic")
			}
		}()
		AsString().Require() // no value is set
	})

	t.Run("AsString().IsValueSet()_is_true", func(t *testing.T) {
		a := AsString().DefaultTo("cat")
		if !a.IsValueSet() {
			t.Errorf("AsString() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsString().IsValueSet()_is_false", func(t *testing.T) {
		a := AsString()
		if a.IsValueSet() {
			t.Errorf("AsString() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsString().SetKey()", func(t *testing.T) {
		e := "NameTest"
		a := AsString().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "good")
		a := AsString().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetValue(cat).Clear().TrySetValue(dog)", func(t *testing.T) {
		e := "dog"
		a := AsString().TrySetValue("cat").Clear().TrySetValue("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsString().TrySetValue(cat).Clear().IsValueSet()", func(t *testing.T) {
		a := AsString().TrySetValue("cat").Clear()
		if a.IsValueSet() {
			t.Errorf("AsString() Failed: expected IsValueSet() to be false")
		}
	})

}

func ExampleAsString() {
	// NOTE: this tests the Print() functionality

	AsString().SetKey("TEST_01").TrySetValue("cat").Print()
	os.Setenv("TEST_VALUE", "dog")
	AsString().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "bear")
	AsString().TrySetByEnv("TEST_VALUE").Print()
	AsString().SetKey("TEST_04").SetValue("bird").PrintMasked()
	AsString().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = cat
	//   TEST_02 = dog
	//   TEST_VALUE = bear
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

func TestAsInt(t *testing.T) {

	t.Run("AsInt().TrySetByString(123)", func(t *testing.T) {
		e := 123
		a := AsInt().TrySetByString("123").Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByString(bad)", func(t *testing.T) {
		e := 456
		a := AsInt().TrySetByString("bad").DefaultTo(456).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByString(cat).Lookup", func(t *testing.T) {
		e := 600
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]int{"cat": 600}
		a := AsInt().TrySetByString("cat").Lookup(table).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().Transform", func(t *testing.T) {
		e := 456
		a := AsInt().Transform(func(chain *IntChain) {
			if !chain.IsValueSet() {
				chain.SetValue(456)
			}
		}).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByEnv(484)", func(t *testing.T) {
		e := 484
		os.Setenv("TEST_VALUE", "484")
		a := AsInt().TrySetByEnv("TEST_VALUE").Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByEnv(bad)", func(t *testing.T) {
		e := 999
		os.Setenv("TEST_VALUE", "bad")
		a := AsInt().TrySetByEnv("TEST_VALUE").DefaultTo(999).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByEnv(cat).Lookup", func(t *testing.T) {
		e := 500
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]int{"cat": 500}
		a := AsInt().TrySetByEnv("TEST_VALUE").Lookup(table).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByEnv(cat).Transform", func(t *testing.T) {
		e := 100
		os.Setenv("TEST_VALUE", "cat")
		a := AsInt().TrySetByEnv("TEST_VALUE").Transform(func(chain *IntChain) {
			val := chain.StringValue()
			if val == "cat" {
				chain.SetValue(100)
			}
		}).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetValue(123).TrySetValue(456)", func(t *testing.T) {
		e := 123
		a := AsInt().TrySetValue(123).TrySetValue(456).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetValue(123).SetValue(456)", func(t *testing.T) {
		e := 456
		a := AsInt().TrySetValue(123).SetValue(456).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt()_null_is_0", func(t *testing.T) {
		e := 0
		a := AsInt().Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsInt() Failed: expected panic")
			}
		}()
		AsInt().Require() // no value is set
	})

	t.Run("AsInt().IsValueSet()_is_true", func(t *testing.T) {
		a := AsInt().DefaultTo(444)
		if !a.IsValueSet() {
			t.Errorf("AsInt() Failed: expected IsValueSet() to be true")
		}
	})

	t.Run("AsInt().IsValueSet()_is_false", func(t *testing.T) {
		a := AsInt()
		if a.IsValueSet() {
			t.Errorf("AsInt() Failed: expected IsValueSet() to be false")
		}
	})

	t.Run("AsInt().SetKey()", func(t *testing.T) {
		e := "NameTest"
		a := AsInt().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsInt() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsInt().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "bad") // NOTE: the name is set even if the value isn't valid
		a := AsInt().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsInt() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsInt().TrySetValue(444).Clear().TrySetValue(887)", func(t *testing.T) {
		e := 887
		a := AsInt().TrySetValue(444).Clear().TrySetValue(887).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetValue(111).Clear().IsValueSet()", func(t *testing.T) {
		a := AsInt().TrySetValue(111).Clear()
		if a.IsValueSet() {
			t.Errorf("AsInt() Failed: expected IsValueSet() to be false")
		}
	})

}

func ExampleAsInt() {
	// NOTE: this tests the Print() functionality

	AsInt().SetKey("TEST_01").TrySetByString("111").Print()
	os.Setenv("TEST_VALUE", "")
	AsInt().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "222")
	AsInt().SetKey("TEST_03").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "333")
	AsInt().TrySetByEnv("TEST_VALUE").Print()
	AsInt().SetKey("TEST_04").SetValue(383).PrintMasked()
	AsInt().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = 111
	//   TEST_02 = 0
	//   TEST_03 = 222
	//   TEST_VALUE = 333
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

func TestAsFloat(t *testing.T) {

	t.Run("AsFloat().TrySetByString(123.25)", func(t *testing.T) {
		e := 123.25
		a := AsFloat().TrySetByString("123.25").Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByString(bad)", func(t *testing.T) {
		e := 456.75
		a := AsFloat().TrySetByString("bad").DefaultTo(456.75).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByString(cat).Lookup", func(t *testing.T) {
		e := 600.5
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]float64{"cat": 600.5}
		a := AsFloat().TrySetByString("cat").Lookup(table).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().Transform", func(t *testing.T) {
		e := 456.75
		a := AsFloat().Transform(func(chain *Float64Chain) {
			if !chain.IsValueSet() {
				chain.SetValue(456.75)
			}
		}).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByEnv(484.9)", func(t *testing.T) {
		e := 484.9
		os.Setenv("TEST_VALUE", "484.9")
		a := AsFloat().TrySetByEnv("TEST_VALUE").Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByEnv(bad)", func(t *testing.T) {
		e := 999.3
		os.Setenv("TEST_VALUE", "bad")
		a := AsFloat().TrySetByEnv("TEST_VALUE").DefaultTo(999.3).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByEnv(cat).Lookup", func(t *testing.T) {
		e := 500.25
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]float64{"cat": 500.25}
		a := AsFloat().TrySetByEnv("TEST_VALUE").Lookup(table).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetByEnv(cat).Transform", func(t *testing.T) {
		e := 100.0
		os.Setenv("TEST_VALUE", "cat")
		a := AsFloat().TrySetByEnv("TEST_VALUE").Transform(func(chain *Float64Chain) {
			val := chain.StringValue()
			if val == "cat" {
				chain.SetValue(100)
			}
		}).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetValue(123).TrySetValue(456)", func(t *testing.T) {
		e := 123.5
		a := AsFloat().TrySetValue(123.5).TrySetValue(456.25).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetValue(123).SetValue(456)", func(t *testing.T) {
		e := 456.75
		a := AsFloat().TrySetValue(123.5).SetValue(456.75).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat()_null_is_0", func(t *testing.T) {
		e := 0.0
		a := AsFloat().Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsFloat() Failed: expected panic")
			}
		}()
		AsFloat().Require() // no value is set
	})

	t.Run("AsFloat().IsSet()_is_true", func(t *testing.T) {
		a := AsFloat().DefaultTo(444)
		if !a.IsValueSet() {
			t.Errorf("AsFloat() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsFloat().IsSet()_is_false", func(t *testing.T) {
		a := AsFloat()
		if a.IsValueSet() {
			t.Errorf("AsFloat() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsFloat().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsFloat().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsFloat().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "bad") // NOTE: the name is set even if the value isn't valid
		a := AsFloat().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsFloat().TrySetValue(444).Clear().TrySetValue(887)", func(t *testing.T) {
		e := 887.25
		a := AsFloat().TrySetValue(444).Clear().TrySetValue(887.25).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetValue(111).Clear().IsSet()", func(t *testing.T) {
		a := AsFloat().TrySetValue(111).Clear()
		if a.IsValueSet() {
			t.Errorf("AsFloat() Failed: expected IsValueSet() to be false")
		}
	})

}

func ExampleAsFloat() {
	// NOTE: this tests the Print() functionality

	AsFloat().SetKey("TEST_01").SetValue(111.25).Print()
	os.Setenv("TEST_VALUE", "222.50")
	AsFloat().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "333.75")
	AsFloat().TrySetByEnv("TEST_VALUE").Print()
	AsFloat().SetKey("TEST_04").SetValue(444.0).PrintMasked()
	AsFloat().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = 111.25
	//   TEST_02 = 222.5
	//   TEST_VALUE = 333.75
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

func TestAsBool(t *testing.T) {

	t.Run("AsBool().TrySetByString(true)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("true").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(bad)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("bad").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(cat).Lookup", func(t *testing.T) {
		e := true
		table := map[string]bool{"cat": true}
		a := AsBool().TrySetByString("cat").Lookup(table).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().Transform", func(t *testing.T) {
		e := true
		a := AsBool().Transform(func(chain *BoolChain) {
			if !chain.IsValueSet() {
				chain.SetValue(true)
			}
		}).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByEnv(true)", func(t *testing.T) {
		e := true
		os.Setenv("TEST_VALUE", "true")
		a := AsBool().TrySetByEnv("TEST_VALUE").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByEnv(bad)", func(t *testing.T) {
		e := true
		os.Setenv("TEST_VALUE", "bad")
		a := AsBool().TrySetByEnv("TEST_VALUE").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByEnv(cat).Lookup", func(t *testing.T) {
		e := true
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]bool{"cat": true}
		a := AsBool().TrySetByEnv("TEST_VALUE").Lookup(table).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByEnv(cat).Transform", func(t *testing.T) {
		e := true
		os.Setenv("TEST_VALUE", "cat")
		a := AsBool().TrySetByEnv("TEST_VALUE").Transform(func(chain *BoolChain) {
			val := chain.StringValue()
			if val == "cat" {
				chain.SetValue(true)
			}
		}).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetValue(false).TrySetValue(true)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetValue(false).TrySetValue(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetValue(false).SetValue(true)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetValue(false).SetValue(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool()_null_is_false", func(t *testing.T) {
		e := false
		a := AsBool().Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsBool() Failed: expected panic")
			}
		}()
		AsBool().Require() // no value is set
	})

	t.Run("AsBool().IsValueSet()_is_true", func(t *testing.T) {
		a := AsBool().DefaultTo(true)
		if !a.IsValueSet() {
			t.Errorf("AsBool() Failed: expected IsValueSet() to be true")
		}
	})

	t.Run("AsBool().IsValueSet()_is_false", func(t *testing.T) {
		a := AsBool()
		if a.IsValueSet() {
			t.Errorf("AsBool() Failed: expected IsValueSet() to be false")
		}
	})

	t.Run("AsBool().SetKey()", func(t *testing.T) {
		e := "NameTest"
		a := AsBool().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsBool() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsBool().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "bad") // NOTE: the name is set even if the value isn't valid
		a := AsBool().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsBool() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsBool().TrySetValue(false).Clear().TrySetValue(true)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetValue(false).Clear().TrySetValue(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetValue(true).Clear().IsValueSet()", func(t *testing.T) {
		a := AsBool().TrySetValue(true).Clear()
		if a.IsValueSet() {
			t.Errorf("AsBool() Failed: expected IsValueSet() to be false")
		}
	})

}

func ExampleAsBool() {
	// NOTE: this tests the Print() functionality

	AsBool().SetKey("TEST_01").TrySetValue(true).Print()
	os.Setenv("TEST_VALUE", "yes")
	AsBool().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "1")
	AsBool().TrySetByEnv("TEST_VALUE").Print()
	AsBool().SetKey("TEST_04").SetValue(false).PrintMasked()
	AsBool().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = true
	//   TEST_02 = true
	//   TEST_VALUE = true
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

func TestAsBoolSpellings(t *testing.T) {

	t.Run("AsBool().TrySetByString(TRUE)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("TRUE").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(TruE)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("TruE").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(yes)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("yes").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(y)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("y").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(1)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetByString("1").Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(FALSE)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("FALSE").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(FaLsE)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("FaLsE").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(no)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("no").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(n)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("n").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetByString(0)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetByString("0").DefaultTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

}

func TestAsSlice(t *testing.T) {

	t.Run("AsSlice().TrySetTo(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSlice().TrySetValue([]string{"cat", "dog"}).Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsString() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().DefaultTo(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSlice().DefaultTo([]string{"cat", "dog"}).Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByString(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSlice().TrySetByString("cat,dog").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByString(cat, dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSlice().TrySetByString("cat, dog").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByString(cat)", func(t *testing.T) {
		e := []string{"cat"}
		a := AsSlice().TrySetByString("cat").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByEnv(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		os.Setenv("TEST_VALUE", "cat,dog")
		a := AsSlice().TrySetByEnv("TEST_VALUE").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByEnv( cat , dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		os.Setenv("TEST_VALUE", " cat , dog")
		a := AsSlice().TrySetByEnv("TEST_VALUE").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByEnv().DefaultTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		os.Setenv("TEST_VALUE", "")
		a := AsSlice().TrySetByEnv("TEST_VALUE").DefaultTo([]string{"dog"}).Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetByEnv(cat).SetValue(dog)", func(t *testing.T) {
		e := []string{"dog"}
		os.Setenv("TEST_VALUE", "cat")
		a := AsSlice().TrySetByEnv("TEST_VALUE").SetValue([]string{"dog"}).Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice().SetValue(cat).SetValue(dog)", func(t *testing.T) {
		e := []string{"dog"}
		a := AsSlice().SetValue([]string{"cat"}).SetValue([]string{"dog"}).Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSlice()_null_is_empty", func(t *testing.T) {
		a := AsSlice().Value()
		if len(a) != 0 || cap(a) != 0 {
			t.Errorf("AsSlice() Failed: expected empty, got \"%v\"", a)
		}
	})

	t.Run("AsSlice()_whitespace_is_empty", func(t *testing.T) {
		a := AsSlice().TrySetByString("  ").Value()
		if len(a) != 0 || cap(a) != 0 {
			t.Errorf("AsSlice() Failed: expected empty, got \"%v\"", a)
		}
	})

	t.Run("AsSlice()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsSlice() Failed: expected panic")
			}
		}()
		AsSlice().Require() // no value is set
	})

	t.Run("AsSlice().IsValueSet()_is_true", func(t *testing.T) {
		a := AsSlice().TrySetByString("cat")
		if !a.IsValueSet() {
			t.Errorf("AsSlice() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsSlice().IsValueSet()_is_false", func(t *testing.T) {
		a := AsSlice()
		if a.IsValueSet() {
			t.Errorf("AsSlice() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsSlice().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsSlice().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsSlice() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsSlice().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "good")
		a := AsSlice().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsSlice() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsSlice().TrySetTo(cat).Clear().TrySetTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		a := AsSlice().TrySetByString("cat").Clear().TrySetByString("dog").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsSlice().TrySetTo(cat).Clear().IsValueSet()", func(t *testing.T) {
		a := AsSlice().TrySetByString("cat").Clear()
		if a.IsValueSet() {
			t.Errorf("AsSlice() Failed: expected IsValueSet() to be false")
		}
	})

	t.Run("AsSlice().UseDelimiter(;).TrySetByString(cat; dog ; bear;)", func(t *testing.T) {
		e := []string{"cat", "dog", "bear"}
		a := AsSlice().UseDelimiter(";").TrySetByString("cat; dog ; bear;").Value()
		if !areSlicesEqual(e, a) {
			t.Errorf("AsSlice() Failed: expected %v, got %v", e, a)
		}
	})

}

func ExampleAsSlice() {
	// NOTE: this tests the Print() functionality

	AsSlice().SetKey("TEST_01").TrySetByString("cat,dog").Print()
	os.Setenv("TEST_VALUE", " cat, blue dog, bear ")
	AsSlice().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "bird")
	AsSlice().TrySetByEnv("TEST_VALUE").Print()
	AsSlice().SetKey("TEST_04").TrySetByString("squirrel").PrintMasked()
	AsSlice().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = [cat dog]
	//   TEST_02 = [cat blue dog bear]
	//   TEST_VALUE = [bird]
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

type Color int

const (
	White Color = iota
	Red
	Yellow
	Blue
	Green
)

func TestAsIntForEnum(t *testing.T) {

	enum := map[string]int{
		"white":  int(White),
		"red":    int(Red),
		"yellow": int(Yellow),
		"blue":   int(Blue),
		"green":  int(Green),
	}

	t.Run("AsInt().TrySetByString(yellow).Lookup()", func(t *testing.T) {
		e := Yellow
		a := AsInt().TrySetByString("yellow").Lookup(enum).Value()
		if Color(a) != e {
			t.Errorf("AsInt()->Enum Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().Lookup().DefaultTo(0)", func(t *testing.T) {
		e := White
		a := AsInt().Lookup(enum).DefaultTo(0).Value()
		if Color(a) != e {
			t.Errorf("AsInt()->Enum Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetByString(bad).Lookup()", func(t *testing.T) {
		e := White
		a := AsInt().TrySetByString("bad").Lookup(enum).Value()
		if Color(a) != e {
			t.Errorf("AsInt()->Enum Failed: expected %d, got %d", e, a)
		}
	})

}

func ExampleAsIntForEnum() {
	// NOTE: this tests the Print() functionality

	enum := map[string]int{
		"white":  int(White),
		"red":    int(Red),
		"yellow": int(Yellow),
		"blue":   int(Blue),
		"green":  int(Green),
	}

	AsInt().SetKey("TEST_01").TrySetByString("yellow").Lookup(enum).PrintLookup(enum)

	// Output:
	//   TEST_01 = yellow
}

func TestAsDuration(t *testing.T) {

	t.Run("AsDuration().TrySetByString(24m)", func(t *testing.T) {
		e := 24 * time.Minute
		a := AsDuration().TrySetByString("24m").Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByString(bad)", func(t *testing.T) {
		e := 0 * time.Second
		a := AsDuration().TrySetByString("bad").Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByString(cat).Lookup", func(t *testing.T) {
		e := 15 * time.Second
		table := map[string]time.Duration{"cat": 15 * time.Second}
		a := AsDuration().TrySetByString("cat").Lookup(table).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().Transform", func(t *testing.T) {
		e := 21 * time.Hour
		a := AsDuration().Transform(func(chain *TimeDurationChain) {
			if !chain.IsValueSet() {
				chain.SetValue(21 * time.Hour)
			}
		}).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByEnv(true)", func(t *testing.T) {
		e := 16*time.Hour + 15*time.Minute
		os.Setenv("TEST_VALUE", "16h15m")
		a := AsDuration().TrySetByEnv("TEST_VALUE").Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByEnv(bad)", func(t *testing.T) {
		e := 13 * time.Hour
		os.Setenv("TEST_VALUE", "bad")
		a := AsDuration().TrySetByEnv("TEST_VALUE").DefaultTo(13 * time.Hour).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByEnv(cat).Lookup", func(t *testing.T) {
		e := 15 * time.Second
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]time.Duration{"cat": 15 * time.Second}
		a := AsDuration().TrySetByEnv("TEST_VALUE").Lookup(table).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetByEnv(cat).Transform", func(t *testing.T) {
		e := 24 * time.Hour
		os.Setenv("TEST_VALUE", "cat")
		a := AsDuration().TrySetByEnv("TEST_VALUE").Transform(func(chain *TimeDurationChain) {
			val := chain.StringValue()
			if val == "cat" {
				chain.SetValue(24 * time.Hour)
			}
		}).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetValue(15m).TrySetValue(17h)", func(t *testing.T) {
		e := 15 * time.Minute
		a := AsDuration().TrySetValue(15 * time.Minute).TrySetValue(17 * time.Hour).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetValue(15m).SetValue(17h)", func(t *testing.T) {
		e := 17 * time.Hour
		a := AsDuration().TrySetValue(15 * time.Minute).SetValue(17 * time.Hour).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration()_null_is_0", func(t *testing.T) {
		e := 0 * time.Second
		a := AsDuration().Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsDuration() Failed: expected panic")
			}
		}()
		AsDuration().Require() // no value is set
	})

	t.Run("AsDuration().IsValueSet()_is_true", func(t *testing.T) {
		a := AsDuration().DefaultTo(15 * time.Minute)
		if !a.IsValueSet() {
			t.Errorf("AsDuration() Failed: expected IsValueSet() to be true")
		}
	})

	t.Run("AsDuration().IsValueSet()_is_false", func(t *testing.T) {
		a := AsDuration()
		if a.IsValueSet() {
			t.Errorf("AsDuration() Failed: expected IsValueSet() to be false")
		}
	})

	t.Run("AsDuration().SetKey()", func(t *testing.T) {
		e := "NameTest"
		a := AsDuration().SetKey("NameTest").Key()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsDuration().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "bad") // NOTE: the name is set even if the value isn't valid
		a := AsDuration().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsDuration().TrySetValue(10s).Clear().TrySetValue(11m)", func(t *testing.T) {
		e := 11 * time.Minute
		a := AsDuration().TrySetValue(10 * time.Second).Clear().TrySetValue(11 * time.Minute).Value()
		if a != e {
			t.Errorf("AsDuration() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsDuration().TrySetValue(true).Clear().IsValueSet()", func(t *testing.T) {
		a := AsDuration().TrySetValue(11 * time.Second).Clear()
		if a.IsValueSet() {
			t.Errorf("AsDuration() Failed: expected IsValueSet() to be false")
		}
	})

}

func ExampleAsDuration() {
	// NOTE: this tests the Print() functionality

	AsDuration().SetKey("TEST_01").TrySetByString("15h13m2s").Print()
	os.Setenv("TEST_VALUE", "13m")
	AsDuration().SetKey("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "17h")
	AsDuration().TrySetByEnv("TEST_VALUE").Print()
	AsDuration().SetKey("TEST_04").SetValue(time.Duration(15 * time.Minute)).PrintMasked()
	AsDuration().SetKey("TEST_05").PrintMasked()

	// Output:
	//   TEST_01 = 15h13m2s
	//   TEST_02 = 13m0s
	//   TEST_VALUE = 17h0m0s
	//   TEST_04 = (set)
	//   TEST_05 = (not-set)
}

func TestIfThenElse(t *testing.T) {

	t.Run("IfThenElse()_then", func(t *testing.T) {
		e := 1
		a := IfThenElse(true, 1, 2)
		if a != e {
			t.Errorf("IfThenElse() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("IfThenElse()_else", func(t *testing.T) {
		e := 2
		a := IfThenElse(false, 1, 2)
		if a != e {
			t.Errorf("IfThenElse() Failed: expected %d, got %d", e, a)
		}
	})

}

/*
func TestResolveAll(t *testing.T) {
	ctx := context.Background()
	list := []*StringChain{
		AsString().Name("SECRET").SetTo("https://pelasne-vaultman.vault.azure.net/secrets/my-secret"),
	}
	ResolveAll(ctx, list).Wait()
	if list[0].Value() != "secret-sauce" {
		t.Errorf("AsInt()->Enum Failed: expected \"%s\", got \"%s\"", "secret-sauce", list[0].Value())
	}
}
*/

// TODO: add tests for Key Vault and AppConfig using mocks
// TODO: add tests for Clamp()
// TODO: add tests for delimiter
// TODO: add tests for Empty()
