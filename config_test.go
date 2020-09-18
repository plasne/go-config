package config

import (
	"context"
	"os"
	"reflect"
	"testing"
)

func TestAsString(t *testing.T) {

	t.Run("AsString().TrySetTo(cat)", func(t *testing.T) {
		e := "cat"
		a := AsString().TrySetTo("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetTo().DefaultTo(cat)", func(t *testing.T) {
		e := "cat"
		a := AsString().TrySetTo("").DefaultTo("cat").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().AllowEmpty().TrySetTo().DefaultTo(cat)", func(t *testing.T) {
		e := ""
		a := AsString().AllowEmpty().TrySetTo("").DefaultTo("cat").Value()
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

	t.Run("AsString().AllowEmpty().TrySetByEnv().DefaultTo(dog)", func(t *testing.T) {
		e := ""
		os.Setenv("TEST_VALUE", "")
		a := AsString().AllowEmpty().TrySetByEnv("TEST_VALUE").DefaultTo("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().TrySetByEnv(cat).SetTo(dog)", func(t *testing.T) {
		e := "dog"
		os.Setenv("TEST_VALUE", "cat")
		a := AsString().TrySetByEnv("TEST_VALUE").SetTo("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString().SetTo(cat).SetTo(dog)", func(t *testing.T) {
		e := "dog"
		a := AsString().SetTo("cat").SetTo("dog").Value()
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

	t.Run("AsString()_whitespace_is_empty", func(t *testing.T) {
		e := ""
		a := AsString().TrySetTo("   ").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsString()_whitespace_is_allowed_on_explict_set", func(t *testing.T) {
		e := "   "
		a := AsString().SetTo("   ").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected \"%s\", got \"%s\"", e, a)
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

	t.Run("AsString().IsSet()_is_true", func(t *testing.T) {
		a := AsString().DefaultTo("cat")
		if !a.IsSet() {
			t.Errorf("AsString() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsString().IsSet()_is_false", func(t *testing.T) {
		a := AsString()
		if a.IsSet() {
			t.Errorf("AsString() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsString().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsString().Name("NameTest").Key()
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

	t.Run("AsString().TrySetTo(cat).Clear().TrySetTo(dog)", func(t *testing.T) {
		e := "dog"
		a := AsString().TrySetTo("cat").Clear().TrySetTo("dog").Value()
		if a != e {
			t.Errorf("AsString() Failed: expected %s, got %s", e, a)
		}
	})

	t.Run("AsString().TrySetTo(cat).Clear().IsSet()", func(t *testing.T) {
		a := AsString().TrySetTo("cat").Clear()
		if a.IsSet() {
			t.Errorf("AsString() Failed: expected IsSet() to be false")
		}
	})

}

func ExampleAsString() {
	// NOTE: this tests the Print() functionality

	AsString().Name("TEST_01").TrySetTo("cat").Print()
	os.Setenv("TEST_VALUE", "dog")
	AsString().Name("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "bear")
	AsString().TrySetByEnv("TEST_VALUE").Print()
	AsString().Name("TEST_04").SetTo("bird").PrintMasked()
	AsString().Name("TEST_05").PrintMasked()

	// Output:
	// TEST_01 = "cat"
	// TEST_02 = "dog"
	// TEST_VALUE = "bear"
	// TEST_04 = (set)
	// TEST_05 = (not-set)
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
			if !chain.IsSet() {
				chain.SetTo(456)
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
				chain.SetTo(100)
			}
		}).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetTo(123).TrySetTo(456)", func(t *testing.T) {
		e := 123
		a := AsInt().TrySetTo(123).TrySetTo(456).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetTo(123).SetTo(456)", func(t *testing.T) {
		e := 456
		a := AsInt().TrySetTo(123).SetTo(456).Value()
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

	t.Run("AsInt().IsSet()_is_true", func(t *testing.T) {
		a := AsInt().DefaultTo(444)
		if !a.IsSet() {
			t.Errorf("AsInt() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsInt().IsSet()_is_false", func(t *testing.T) {
		a := AsInt()
		if a.IsSet() {
			t.Errorf("AsInt() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsInt().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsInt().Name("NameTest").Key()
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

	t.Run("AsInt().TrySetTo(444).Clear().TrySetTo(887)", func(t *testing.T) {
		e := 887
		a := AsInt().TrySetTo(444).Clear().TrySetTo(887).Value()
		if a != e {
			t.Errorf("AsInt() Failed: expected %d, got %d", e, a)
		}
	})

	t.Run("AsInt().TrySetTo(111).Clear().IsSet()", func(t *testing.T) {
		a := AsInt().TrySetTo(111).Clear()
		if a.IsSet() {
			t.Errorf("AsInt() Failed: expected IsSet() to be false")
		}
	})

}

func ExampleAsInt() {
	// NOTE: this tests the Print() functionality

	AsInt().Name("TEST_01").TrySetByString("111").Print()
	os.Setenv("TEST_VALUE", "")
	AsInt().Name("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "222")
	AsInt().Name("TEST_03").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "333")
	AsInt().TrySetByEnv("TEST_VALUE").Print()
	AsInt().Name("TEST_04").SetTo(383).PrintMasked()
	AsInt().Name("TEST_05").PrintMasked()

	// Output:
	// TEST_01 = 111
	// TEST_02 = 0
	// TEST_03 = 222
	// TEST_VALUE = 333
	// TEST_04 = (set)
	// TEST_05 = (not-set)
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
		a := AsFloat().Transform(func(chain *FloatChain) {
			if !chain.IsSet() {
				chain.SetTo(456.75)
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
		a := AsFloat().TrySetByEnv("TEST_VALUE").Transform(func(chain *FloatChain) {
			val := chain.StringValue()
			if val == "cat" {
				chain.SetTo(100)
			}
		}).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetTo(123).TrySetTo(456)", func(t *testing.T) {
		e := 123.5
		a := AsFloat().TrySetTo(123.5).TrySetTo(456.25).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetTo(123).SetTo(456)", func(t *testing.T) {
		e := 456.75
		a := AsFloat().TrySetTo(123.5).SetTo(456.75).Value()
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
		if !a.IsSet() {
			t.Errorf("AsFloat() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsFloat().IsSet()_is_false", func(t *testing.T) {
		a := AsFloat()
		if a.IsSet() {
			t.Errorf("AsFloat() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsFloat().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsFloat().Name("NameTest").Key()
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

	t.Run("AsFloat().TrySetTo(444).Clear().TrySetTo(887)", func(t *testing.T) {
		e := 887.25
		a := AsFloat().TrySetTo(444).Clear().TrySetTo(887.25).Value()
		if a != e {
			t.Errorf("AsFloat() Failed: expected %f, got %f", e, a)
		}
	})

	t.Run("AsFloat().TrySetTo(111).Clear().IsSet()", func(t *testing.T) {
		a := AsFloat().TrySetTo(111).Clear()
		if a.IsSet() {
			t.Errorf("AsFloat() Failed: expected IsSet() to be false")
		}
	})

}

func ExampleAsFloat() {
	// NOTE: this tests the Print() functionality

	AsFloat().Name("TEST_01").TrySetTo(111.25).Print()
	os.Setenv("TEST_VALUE", "222.50")
	AsFloat().Name("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "333.75")
	AsFloat().TrySetByEnv("TEST_VALUE").Print()
	AsFloat().Name("TEST_04").SetTo(444.0).PrintMasked()
	AsFloat().Name("TEST_05").PrintMasked()

	// Output:
	// TEST_01 = 111.250000
	// TEST_02 = 222.500000
	// TEST_VALUE = 333.750000
	// TEST_04 = (set)
	// TEST_05 = (not-set)
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
		os.Setenv("TEST_VALUE", "cat")
		table := map[string]bool{"cat": true}
		a := AsBool().TrySetByString("cat").Lookup(table).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().Transform", func(t *testing.T) {
		e := true
		a := AsBool().Transform(func(chain *BoolChain) {
			if !chain.IsSet() {
				chain.SetTo(true)
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
				chain.SetTo(true)
			}
		}).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetTo(false).TrySetTo(true)", func(t *testing.T) {
		e := false
		a := AsBool().TrySetTo(false).TrySetTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetTo(false).SetTo(true)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetTo(false).SetTo(true).Value()
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

	t.Run("AsBool().IsSet()_is_true", func(t *testing.T) {
		a := AsBool().DefaultTo(true)
		if !a.IsSet() {
			t.Errorf("AsBool() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsBool().IsSet()_is_false", func(t *testing.T) {
		a := AsBool()
		if a.IsSet() {
			t.Errorf("AsBool() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsBool().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsBool().Name("NameTest").Key()
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

	t.Run("AsBool().TrySetTo(false).Clear().TrySetTo(true)", func(t *testing.T) {
		e := true
		a := AsBool().TrySetTo(false).Clear().TrySetTo(true).Value()
		if a != e {
			t.Errorf("AsBool() Failed: expected %t, got %t", e, a)
		}
	})

	t.Run("AsBool().TrySetTo(true).Clear().IsSet()", func(t *testing.T) {
		a := AsBool().TrySetTo(true).Clear()
		if a.IsSet() {
			t.Errorf("AsBool() Failed: expected IsSet() to be false")
		}
	})

}

func ExampleAsBool() {
	// NOTE: this tests the Print() functionality

	AsBool().Name("TEST_01").TrySetTo(true).Print()
	os.Setenv("TEST_VALUE", "yes")
	AsBool().Name("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "1")
	AsBool().TrySetByEnv("TEST_VALUE").Print()
	AsBool().Name("TEST_04").SetTo(false).PrintMasked()
	AsBool().Name("TEST_05").PrintMasked()

	// Output:
	// TEST_01 = true
	// TEST_02 = true
	// TEST_VALUE = true
	// TEST_04 = (set)
	// TEST_05 = (not-set)
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

func TestAsSplice(t *testing.T) {

	t.Run("AsSplice().TrySetTo(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSplice().TrySetTo([]string{"cat", "dog"}).Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsString() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().DefaultTo(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSplice().DefaultTo([]string{"cat", "dog"}).Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByString(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSplice().TrySetByString("cat,dog").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByString(cat, dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		a := AsSplice().TrySetByString("cat, dog").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByString(cat)", func(t *testing.T) {
		e := []string{"cat"}
		a := AsSplice().TrySetByString("cat").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByEnv(cat,dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		os.Setenv("TEST_VALUE", "cat,dog")
		a := AsSplice().TrySetByEnv("TEST_VALUE").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByEnv( cat , dog)", func(t *testing.T) {
		e := []string{"cat", "dog"}
		os.Setenv("TEST_VALUE", " cat , dog")
		a := AsSplice().TrySetByEnv("TEST_VALUE").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByEnv().DefaultTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		os.Setenv("TEST_VALUE", "")
		a := AsSplice().TrySetByEnv("TEST_VALUE").DefaultTo([]string{"dog"}).Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetByEnv(cat).SetTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		os.Setenv("TEST_VALUE", "cat")
		a := AsSplice().TrySetByEnv("TEST_VALUE").SetTo([]string{"dog"}).Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice().SetTo(cat).SetTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		a := AsSplice().SetTo([]string{"cat"}).SetTo([]string{"dog"}).Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected \"%v\", got \"%v\"", e, a)
		}
	})

	t.Run("AsSplice()_null_is_empty", func(t *testing.T) {
		a := AsSplice().Value()
		if len(a) != 0 || cap(a) != 0 {
			t.Errorf("AsSplice() Failed: expected empty, got \"%v\"", a)
		}
	})

	t.Run("AsSplice()_whitespace_is_empty", func(t *testing.T) {
		a := AsSplice().TrySetByString("  ").Value()
		if len(a) != 0 || cap(a) != 0 {
			t.Errorf("AsSplice() Failed: expected empty, got \"%v\"", a)
		}
	})

	t.Run("AsSplice()_required_panics", func(t *testing.T) {
		defer func() {
			if err := recover(); err == nil {
				t.Error("AsSplice() Failed: expected panic")
			}
		}()
		AsSplice().Require() // no value is set
	})

	t.Run("AsSplice().IsSet()_is_true", func(t *testing.T) {
		a := AsSplice().TrySetByString("cat")
		if !a.IsSet() {
			t.Errorf("AsSplice() Failed: expected IsSet() to be true")
		}
	})

	t.Run("AsSplice().IsSet()_is_false", func(t *testing.T) {
		a := AsSplice()
		if a.IsSet() {
			t.Errorf("AsSplice() Failed: expected IsSet() to be false")
		}
	})

	t.Run("AsSplice().Key()_from_Name()", func(t *testing.T) {
		e := "NameTest"
		a := AsSplice().Name("NameTest").Key()
		if a != e {
			t.Errorf("AsSplice() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsSplice().Key()_from_TrySetByEnv()", func(t *testing.T) {
		e := "TEST_VALUE"
		os.Setenv("TEST_VALUE", "good")
		a := AsSplice().TrySetByEnv("TEST_VALUE").Key()
		if a != e {
			t.Errorf("AsSplice() Failed: expected \"%s\", got \"%s\"", e, a)
		}
	})

	t.Run("AsSplice().TrySetTo(cat).Clear().TrySetTo(dog)", func(t *testing.T) {
		e := []string{"dog"}
		a := AsSplice().TrySetByString("cat").Clear().TrySetByString("dog").Value()
		if !reflect.DeepEqual(e, a) {
			t.Errorf("AsSplice() Failed: expected %v, got %v", e, a)
		}
	})

	t.Run("AsSplice().TrySetTo(cat).Clear().IsSet()", func(t *testing.T) {
		a := AsSplice().TrySetByString("cat").Clear()
		if a.IsSet() {
			t.Errorf("AsSplice() Failed: expected IsSet() to be false")
		}
	})

}

func ExampleAsSplice() {
	// NOTE: this tests the Print() functionality

	AsSplice().Name("TEST_01").TrySetByString("cat,dog").Print()
	os.Setenv("TEST_VALUE", " cat, blue dog, bear ")
	AsSplice().Name("TEST_02").TrySetByEnv("TEST_VALUE").Print()
	os.Setenv("TEST_VALUE", "bird")
	AsSplice().TrySetByEnv("TEST_VALUE").Print()
	AsSplice().Name("TEST_04").TrySetByString("squirrel").PrintMasked()
	AsSplice().Name("TEST_05").PrintMasked()

	// Output:
	// TEST_01 = [cat dog]
	// TEST_02 = [cat blue dog bear]
	// TEST_VALUE = [bird]
	// TEST_04 = (set)
	// TEST_05 = (not-set)
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

	AsInt().Name("TEST_01").TrySetByString("yellow").Lookup(enum).PrintLookup(enum)

	// Output:
	// TEST_01 = yellow
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

// TODO: add tests for Key Vault and AppConfig
