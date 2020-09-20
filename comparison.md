<!-- markdownlint-disable no-hard-tabs -->

# Comparison to existing modules

I looked at a number of modules to address configuration management via environment variable. I spent the most time with <https://github.com/sethvargo/go-envconfig> so I will discuss that here. Please don't take this as a criticism of this module, it seems very thoughtfully constructed and can easily address a lot of scenarios. However, I created this module because I think I can fix a few things and address some additional scenarios.

## Scenarios I believe are handled incorrectly

I have illustrated 2 problems I found with the module below in unit tests.

* If a variable is not assigned a value, the default is used. However, if a value is assigned, but that value is not appropriate for the listed datatype, then the default value of the datatype is used, not the specified default. I think there are 2 logical ways to handle this scenario, either the module should use the provided default or should panic(). In my module, I use the provided default.

* If a provided value is effectively empty, like passing "" to a splice, then go-envconfig does ignore the value, but it returns an array of cap=0,len=0 instead of the provided default. In my module, I use the provided default.

```go
func TestGoEnvConfig(t *testing.T) {
	ctx := context.Background()

	type Config struct {
		Int          int      `env:"INT"`
		IntDefault   int      `env:"INT,default=7"`
		Array        []string `env:"ARRAY"`
		ArrayDefault []string `env:"ARRAY,default=cat,dog"`
	}

	// NOTE: I believe this should be 7, not 0 because "dog" is not an integer
	t.Run("Int_non_numeric", func(t *testing.T) {
		var config Config
		lookuper := envconfig.MapLookuper(map[string]string{
			"INT": "dog",
		})
		if err := envconfig.ProcessWith(ctx, &config, lookuper); err != nil {
			t.Errorf("could not get configuration")
		}
		e := 7
		if !reflect.DeepEqual(config.IntDefault, e) {
			t.Errorf("int Failed: expected %v, got %v", e, config.IntDefault)
		}
	})

	// NOTE: I believe this should be [cat dog] because the array does not have a value, instead it returns cap:0,len:0
	t.Run("Array_default_from_empty", func(t *testing.T) {
		var config Config
		lookuper := envconfig.MapLookuper(map[string]string{
			"ARRAY": "",
		})
		if err := envconfig.ProcessWith(ctx, &config, lookuper); err != nil {
			t.Errorf("could not get configuration")
		}
		e := []string{"cat", "dog"}
		if !reflect.DeepEqual(config.ArrayDefault, e) {
			t.Errorf("array Failed: expected %v, got %v", e, config.ArrayDefault)
		}
	})

}
```

```text
--- FAIL: TestGoEnvConfigInt (0.00s)
    --- FAIL: TestGoEnvConfigInt/Int_non_numeric (0.00s)
        config_test.go:32: could not get configuration
        config_test.go:36: int Failed: expected 7, got 0
    --- FAIL: TestGoEnvConfigInt/Array_default_from_empty (0.00s)
        config_test.go:64: could not get configuration
        config_test.go:68: array Failed: expected [cat dog], got []
```

## Additional scenarios I wanted to support

* DotEnv - Loading environment variables from a file can be really useful for local development. This module supports dotenv via <https://github.com/joho/godotenv>.

* Azure AppConfig - If the application has a lot of variables, or you want to have a single place for all your configurations, it can be helpful to use a configuration datastore. This module supports AppConfig for this purpose.

* Azure Key Vault - Often you do not want secrets stored in clear text, so this module supports using Azure Key Vault to resolve a URL in env. It also supports getting those URLs from Azure AppConfig.

* Multiple Environment Variables - A common scenario is to have a generic variable, like DB_CONNSTRING, and then have more specific variables like HISTORY_DB_CONNSTRING. The configuration should use HISTORY_DB_CONNSTRING if it is provided, but then fall back to DB_CONNSTRING if it isn't provided. This module supports checking as many environment variables as you want until you get a legitimate value.

* Lookup - Generally a lookup can be used for Enums that are based off an integer of some kind. This scenario is supported by this module.

* Clamp - Often numeric values need to be within some range, this module allows you to clamp the value between a min and max value.

* Transform - For more exotic cases, this module supports a transform operation that allows you to write code to modify the current value.

* Print - I believe that all configuration values should be printed to the console with the application starts up so it is clear to the user what configuration is being used. To facilitate this, the module supports printing each key/value pair. It also supports printing with the value masked or showing a label instead of an int.

Some of these scenarios could be supported using go-envconfig using extensions <https://github.com/sethvargo/go-envconfig#extension>, but there are several reasons this is not sufficient:

* You would have to embed code for the extension in every project or wrap them into a module.

* The extensions work on the raw values before the destination data type is determined. This creates the odd scenario where you are parsing a value and potentially changing before it is being parsed again by the module.

* The module supports mapping the env to multiple variables. While I doubt this is a common scenario it procludes you from using different extension for different variables.
