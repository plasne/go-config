<!-- markdownlint-disable no-hard-tabs -->

# go-config

Every application has to handle common scenarios like configuration management. I wrote this Go module to address configuration via environment variables because I did not find any existing modules that handled all conditions in the way I thought they should be handled or addressed all the scenarios I had in mind (read more [here](./comparison.md)).

## Why environment variables over flags

Generally I use environment variables for configuration of services (things running on a server somewhere) and flags for configuration of tools (command line tools that a user runs).

Many platforms have a native way of managing environment variables. For example, Kubernetes allows them to be defined in the manifest and even supports setting secrets from a vault, Azure App Service sets configuration settings as environment variables, a Dockerfile can contain ENV and supports override, and so on.

I don't like hardcoding values. If there is a some choice to be made about what to hardcode a variable as (ex. concurrency, thresholds, etc.), I make an environment variable for that option and default it. This allows me to easily performance tune the application in different environments later. Doing this means a robust service could have dozens or even hundreds of configuration points, but you may only have to set a very small number to have a working solution.

While developing an application locally, it can be very handy to put a lot of configuration settings into a .env file (supported out-of-the-box) and then .gitignore that file. This gives you an easy way to work with a lot of configuration options during development (even secrets) without checking them in. You can just put an underscore in front of specific variables to disable them. Or you could have multiple .env files with different configurations for testing.

:warning: A colleague raised a concern about environment variables being buried throughout the code introducing an element of "magic" because it may be difficult to determine what settings are available and how they are set. I agree, I would recommend 3 mitigations:

* Always set your entire configuration at startup (ie. all envs are set in one place) - I like to use init().

* Always put your entire configuration into a configuration struct (ie. there are no mystery settings).

* Always print the configuration at startup (ie. there is no confusion on what configuration was used or how variables were interrogate and transformed).

The complete sample at the bottom of this guide shows an implementation of this.

## Installation

```bash
go get github.com/plasne/go-config
```

## Usage

In init(), you want to call Startup() and then set your variables by chaining all the necessarily methods. For example, a simple case of setting a single string value from an environment variable named SAMPLE and defaulting to "cat" if the variable was not set...

```go
package main

import (
	goconfig "github.com/plasne/go-config"
)

func main() {

	// startup
	err := goconfig.Startup()
	if err != nil {
		panic(err)
	}

	var myString := goconfig.AsString().TrySetByEnv("SAMPLE").DefaultTo("cat").Value()

}
```

## Build

If you make any changes to AsDataType.go, you must generate gen-AsDataType.go again using <https://github.com/cheekybits/genny>.

To generate, you need to...

```bash
go generate
```

## Common Scenarios

This shows some common scenarios to get you started. The next section details all options and a complete sample is available at the bottom of this page.

```go
// SCENARIO: set a required value
// pull a string value from an environment variable, print it, panic() if not supplied, set the variable to the value
STORAGE_ACCOUNT := goconfig.AsString().TrySetByEnv("STORAGE_ACCOUNT").Print().Require().Value()

// SCENARIO: set a password
// pull a string value from an environment variable, potentially resolve it in Key Vault, print whether or not it was set, panic() if not supplied, set the variable to the value
ctx := context.Background()
STORAGE_KEY := goconfig.AsString().TrySetByEnv("STORAGE_KEY").Resolve(ctx).PrintMasked().Require().Value()

// SCENARIO: set a reasonable numeric value
// pull a string value from an environment variable, parse it into an int if possible, if not set - take the default, clamp the value between 1 and 256, print it, set the variable to the value
CONCURRENCY := goconfig.AsInt().TrySetByEnv("CONCURRENCY").DefaultTo(8).Clamp(1, 256).Print().Value()

// SCENARIO: choose the most specific value supplied from multiple options
// pull a string value from the HISTORY_DB_CONNSTRING env (the more specific setting), if that was empty - pull a string value from DB_CONNSTRING env (the more generic setting), print it, panic() if neither were non-empty, set the variable to the value
HISTORY_DB_CONNSTRING := goconfig.AsString().TrySetByEnv("HISTORY_DB_CONNSTRING").TrySetByEnv("DB_CONNSTRING").Print().Require().Value()

// SCENARIO: set an enum that is incremented as an int
// pull a string value from an environment variable, it is probably one of the supported strings (but if it were a number - that would be fine), lookup the string to convert to an int, clamp it to a supported value, default it to 0 if not provided, print the label, set the variable to the value
table := map[string]int{
	"env": int(AuthMode_Env),
	"cli": int(AuthMode_Cli),
}
GOCONFIG_AUTH_MODE := authMode(AsInt().TrySetByEnv("GOCONFIG_AUTH_MODE").Lookup(table).Clamp(0, 1).DefaultTo(0).PrintLookup(table).Value())

// SCENARIO: allow a flag to override the env
// try to set based on the flag (if it is 0 - it won't be set), try to set by env var, default to 8 if neither worked, clamp between 1 and 256, print the value, set the variable to the value
var concurrency int
flag.IntVar(&concurrency, "concurrency", 0, "Sets the number of calls made in parallel for table operations.")
flag.Parse()
CONCURRENCY := goconfig.AsInt().TrySetValue(concurrency).TrySetByEnv("CONCURRENCY").DefaultTo(8).Clamp(1, 256).Print().Value()

// SCENARIO: transform a value that may not be in the right format
// pull a string value from an environment variable, if the value is set - transform it into a URL if it isn't already, print it, set the variable to the value
GOCONFIG_APPCONFIG := AsString().TrySetByEnv("GOCONFIG_APPCONFIG").Transform(func(chain *StringChain) {
	if chain.IsValueSet() {
		val := strings.ToLower(chain.Value())
		if !strings.HasPrefix(val, "https://") {
			val = "https://" + val
		}
		if strings.HasSuffix(val, "/") {
			val = strings.TrimRight(val, "/")
		}
		if !strings.HasSuffix(val, ".azconfig.io") {
			val += ".azconfig.io"
		}
		chain.SetValue(val)
	}
}).Print().Value()

// SCENARIO: only accept a URL
// pull a string value from an environment variable, clear it if it isn't a URL, print it, set the variable to the value
URL := AsString().TrySetByEnv("URL").Transform(func(chain *StringChain) {
	if chain.IsValueSet() {
		val := strings.ToLower(chain.Value())
		if !strings.HasPrefix(val, "http://") && !strings.HasPrefix(val, "https://") {
			chain.Clear()
		}
	}
}).Print().Value()

// SCENARIO: parse an int from a string obtained some other way
// set a key so print knows how to show it, try to parse a string into an int, print it, set the variable to the value
// OUTPUT: VALUE = 17
VALUE := AsInt().SetKey("VALUE").TrySetByString("17").Print().Value()

// SCENARIO: parse a string into a slice
// set a key so print knows how to show it, try to parse a string into a slice, print it, set the variable to the value
// OUTPUT: SLICE = [dog cat bear]
SLICE := AsSlice().SetKey("SLICE").UseDelimiter(";").TrySetByString("dog; cat; bear;").Print().Value()
```

## Datatypes and Methods

The following datatypes are supported:

| Method | Golang datatype | Empty | Methods | Notes |
| ---- | ---- | ---- | ---- | ---- |
| AsInt() | int | 0 | Offers Clamp(), PrintLookup(). | |
| AsFloat() | float64 | 0.0 | Offers Clamp(). | |
| AsString() | string | "" | | strval and value are always the same. |
| AsBool() | bool | false | | Supports true, yes, y, or 1 for TRUE. Supports false, no, n, or 0 for FALSE. |
| AsDuration() | time.Duration | time.Duration(0) | | |
| AsSlice() | []string | []string{} cap=0, len=0 | Offers UseDelimiter(). | Delimited on comma by default. Whitespace is trimmed from the left and right of each entry. |

Before we get into the chain, there is an important concept. All of the datatypes have storage for name, strval, and value. The strval is set the first time a non-empty string is provided. The value is set the first time a string is provided that can be parsed successfully or a method provides a value in the datatype natively. The strval is useful for methods like Lookup() and Transform() that might want to deal with some kind of label that will be translated into a value of the appropriate datatype. The name is only used for Print() to show a meaningful key/value pair.

Each of the datatype methods will start a chain that allows for any number of the following:

* __SetKey(name string)__ - You supply a name for the chain which will show as the key in the key/value pair that is shown by Print(). If you aren't going to Print(), you don't need to specify a key.

* __SetStringValue(value string)__ - You supply a string and this method sets strval.

* __SetValue(value datatype)__ - You supply a value in the appropriate datatype. The value is set to the provided value regardless of whether or not it has been previously set.

* __SetEmpty(value datatype)__ - You supply a value in the appropriate datatype. When you call any of the Try-prefixed methods, this module only sets the value if the value provided is not empty. The SetEmpty() method allows you to change the definition of empty from what is shown in the datatype table above. For example, you might want an int's default to be -1 if 0 is a legitimate value.

* __Clear()__ - This clears the value if it is set. It has no impact on name or strval. The primary purpose of this method is to revoke a set value in a Transform().

* __TrySetValue(value datatype)__ - You supply a value in the appropriate datatype. If a value has not been set yet, it will be set to this value.

* __DefaultTo(value datatype)__ - This is simply an alias for TrySetTo().

* __TrySetByEnv(name string)__ - You supply the name of an environment variable. If there is not a name specified, name will be set as provided to this method. This method will read the environment variable of the specified name as a string and store it as strval provided the string is non-empty and strval has not set already. If a value has not been set yet, it will then attempt to parse the string to the specified datatype. If successful, the value will be set. To clarify, this method attempts to set the name, strval, and value independently.

* __TrySetByString(value string)__ - You supply a string value. If a strval has not been set yet, this method will store it as strval provided the string is non-empty. If a value has not been set yet, it will then attempt to parse the string to the specified datatype. If successful, the value will be set. To clarify, this method attempts to set the strval and value independently. AsString() does not have this method, you can use TrySetTo() instead.

* __Lookup(map[string]datatype)__ - You specify a map. The strval (or value for AsString()) is used as the key to return a value of the specified datatype. The chain has its value set to the found value even if it was previously set. If strval was not set or a match was not found, this method changes nothing. The key is tried with its provided casing and as all lowercase.

* __Transform(func (*chain) {})__ - You provide a func() that can use any methods in the chain with any logic you want to set, clear, transform, etc. any values.

* __Resolve(ctx context.Context)__ - You provide a context and if the strval (or value for AsString()) is an Azure Key Vault Secret URL the secret will be read from Key Vault. Provided it can be parsed into the correct datatype, it will be set as the value even if a value was previously set. If the strval (or value for AsString()) was not set or was not an Azure Key Vault Secret URL, this method changes nothing.

* __Print()__ - The Key() and Value() methods are called and then printed to the console as "key = value".

* __PrintMasked()__ - The Key() method is called and then printed to the console as "key = (set)" or "key = (not-set)" depending on whether or not a value has been set.

* __PrintLookup(map[string]int)__ - This is only available on AsInt(). You supply a map (typically the same as you might have supplied to Lookup()) and "key = lookup" will be printed. In other words, rather than printing a numeric value, you can print a label.

* __Require()__ - This panics if the value is not set.

* __Clamp(min datatype, max datatype)__ - This is only available on numeric types (AsInt() and AsFloat()). You supply a minimum and maximum value and if the value is set, it is fixed inside this range.

* __UseDelimiter(delimiter string)__ - This is only available on AsSlice(). You supply a delimiter to use instead of comma to separate a provided string into a slice.

The chain can be completed with any of these (but they do not continue the chain):

* __Key()__ - This returns the name of the variable. This will have come from Name() or the first call to TrySetByEnv().

* __Value()__ - This returns the value in the specified datatype. If a value is not set, an empty value will be returned.

* __StringValue()__ - This returns the strval if it is set or an empty string if it wasn't. This is most commonly used in Transform().

* __IsValueSet()__ - This returns true or false depending on whether the value is set for the specified datatype. This is most commonly used in Transform().

* __IsStringValueSet()__ - This returns true or false depending on whether the strval is set. This is most commonly used in Transform().

## Startup()

The Startup() method does the following:

1. Looks for a .env file and processes it if present.

2. Resolves and prints the pre-configuration variables (GOCONFIG_AUTH_MODE, GOCONFIG_APPCONFIG, and GOCONFIG_CONFIG_KEYS).

3. Loads environment variables from AppConfig if appropriate.

## DotEnv

<https://github.com/joho/godotenv> is already referenced in Startup(), so the module will read a .env file without any additional configuration. See the documentation at that link for more details.

## AppConfig

To support AppConfig, you must specify the following environment variables:

* GOCONFIG_CREDS (default: "default") - This is a comma-delimited list of credential types to support. This can be set to any of the following: "default" (`DefaultAzureCredential`), "env" (`EnvironmentCredential`), "mi" (`ManagedIdentityCredential`), or "cli" (`AzureCLICredential`). You can use a

This can be set to "env" or "cli". Either way, this leverages the Azure Go SDK to authenticate the REST calls to AppConfig and/or Key Vault. When set to "env", you can authenticate via Client Credentials, Client Certificate, Resource Owner Password, or Azure Managed Service Identity depending on how you configure additional environment variables. When set to "cli", provided you have az-cli installed, it will authenticate using your current credentials (you may need to run "az login" first). You can get more details here: <https://github.com/Azure/azure-sdk-for-go#more-authentication-details>.

:information_source: Without setting GOCONFIG_AUTH_MODE or any other environment variables, the solution will attempt to use the local MSI endpoint.

* GOCONFIG_APPCONFIG (REQUIRED) - You must specify the name or the full URL to your Azure AppConfig instance (ex. <https://pelasne-config.azconfig.io>).

* GOCONFIG_CONFIG_KEYS (REQUIRED) - You must provide a comma-separated list of key filters. All key/value pairs that match the filters will be considered. Filters are applied from left to right and if a key already exists, it will be ignored. The "key" used will be last colon-separated section of the key. You can find out more about key filters here: <https://github.com/Azure/AppConfiguration/blob/main/docs/REST/kv.md#filtering>.

:warning: The account that is used for authentication must have the "App Configuration Data Reader" role even if it has "Contributor" or "Owner". Also note that it can take up to 30 minutes for this new role to take effect. You will get an HTTP 403 if this role is not provided.

Consider the following example of values stored in AppConfig (exported from AppConfig)...

```json
{
    "override:CONCURRENCY": "8",
    "sample:CONCURRENCY": "32",
    "sample:SERVER_HOST_URL": "http://auth.plasne.com"
}
```

You could get these using the following filters...

```text
GOCONFIG_CONFIG_KEYS=override:*, sample:*
```

All 3 values would be fetched and evaluated into 2 separate key/value pairs: CONCURRENCY=8 and SERVER_HOST_URL=<http://auth.plasne.com>. Notice that CONCURRENCY was provided twice, but the "override" prefix took precident because the rules are evaluated from left to right. This allows you to build a standard configuration and then apply higher precident configuration items. You are not required to use colon-separated keys, but it allows you to implement this pattern easily.

You can make the keys as complicated as you like, for instance I often use "instance:service:environment:key".

AppConfig supports storing Key Vault URLs for secrets, this is fully supported and the URL will be extracted and can work with the Resolve() method.

:warning: Pulling key/value pairs from AppConfig can take a while on a cold start. It is common that it might take 60-90 seconds.

## Key Vault

To support Key Vault, you must specify GOCONFIG_AUTH_MODE as described above.

If the Resolve() method is called and a legitimate Azure Key Vault Secret URL is currently the string value of the variable, the secret will be fetched and made the new string value.

## Complete Sample

Below is a sample of normal usage. Take particular note of a few things:

* The startup will load a .env file (if there is one), display all pre-configuration settings, pull CONFIG_KEYS from APPCONFIG (if those are specified), and then set those as environment variables.

* Setting Enums using AsInt is supported as shown below for the log levels.

* All variables Print() after they are processed so the user can look at the logs to see what happened.

* The STORAGE_KEY is a secret so it may be stored in Key Vault. The Resolve() method will fetch it if it is, and it won't hurt anything if it isn't. You would not want to print the actual value to the console, but you will still want to see that it was set by using PrintMasked().

* Variables that call Require() will panic if the value is not set successfully.

In addition, I personally like the following patterns:

* Loading all configuration variables in init() instead of main().

* Putting the configuration items into a struct.

* Setting the configuration as a local variable scoped for the entire module.

* Using all uppercase (even though they aren't technically constants, they generally should be treated the same) where the name matches the environment variable names.

```go
import (
	"goconfig" github.com/plasne/go-config
)

type Config struct {
	STORAGE_ACCOUNT string
	STORAGE_KEY     string
	RETENTION       time.Duration
	CONCURRENCY     int
	INTERVAL        time.Duration
}

var config Config

func init() {

	// startup config
	err := goconfig.Startup(ctx)
	if err != nil {
		panic(err)
	}

	// start config block
	fmt.Println("CONFIGURATION:")

	// configure logging
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	logLevels := map[string]int{
		"trace":    int(zerolog.TraceLevel),
		"debug":    int(zerolog.DebugLevel),
		"info":     int(zerolog.InfoLevel),
		"warn":     int(zerolog.WarnLevel),
		"error":    int(zerolog.ErrorLevel),
		"fatal":    int(zerolog.FatalLevel),
		"panic":    int(zerolog.PanicLevel),
		"nolevel":  int(zerolog.NoLevel),
		"disabled": int(zerolog.Disabled),
	}
	logLevel := goconfig.AsInt().TrySetByEnv("LOG_LEVEL").Lookup(logLevels).DefaultTo(int(zerolog.InfoLevel)).PrintLookup(logLevels).Value()
	zerolog.SetGlobalLevel(zerolog.Level(logLevel))

	// load configuration
	config.STORAGE_ACCOUNT = goconfig.AsString().TrySetByEnv("STORAGE_ACCOUNT").Print().Require().Value()
	config.STORAGE_KEY = goconfig.AsString().TrySetByEnv("STORAGE_KEY").Resolve(ctx).PrintMasked().Require().Value()
	config.RETENTION = goconfig.AsDuration().TrySetByEnv("RETENTION").DefaultTo(24 * time.Hour).Print().Value()
	config.CONCURRENCY = goconfig.AsInt().TrySetByEnv("CONCURRENCY").DefaultTo(8).Clamp(1, 256).Print().Value()
	config.INTERVAL = goconfig.AsDuration().TrySetByEnv("INTERVAL").DefaultTo(10 * time.Second).Print().Value()

}
```

Given a .env file like this...

```text
AUTH_MODE=cli
APPCONFIG=pelasne-config
CONFIG_KEYS=override:*, sample:*
CONCURRENCY=32
STORAGE_ACCOUNT=pelasnediagdiag
STORAGE_KEY=W...Q==
HOURS_TO_RETAIN=6
```

...the output will look something like this...

```text
PRE-CONFIGURATION:
  GOCONFIG_AUTH_MODE = cli
  GOCONFIG_APPCONFIG = "https://pelasne-config.azconfig.io"
  GOCONFIG_CONFIG_KEYS = [override:* sample:*]
CONFIGURATION:
  LOG_LEVEL = info
  STORAGE_ACCOUNT = "pelasnediagdiag"
  STORAGE_KEY = (set)
  RETENTION = 6h0m0s
  CONCURRENCY = 32
  INTERVAL = 10s
```

## Future

There are some things I would like to expand in the future, including...

* Finish unit tests including mocks.

* "Clamp" for strings, ie. some way to ensure the value is within a specific list.

* Support loading .env which allows you to specify another .env file to load.
