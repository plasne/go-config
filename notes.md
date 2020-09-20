<!-- markdownlint-disable no-hard-tabs -->

# Notes

I used this project to learn to code in Go. As such, there were some paths I went down that I didn't pursue, but I wanted to capture the code for use in later projects. I will probably move this somewhere more permanent, but for now...

## Getting an Access Token from the az-cli cache file

It turns out this is supported by the Azure Go SDK anyway, so it wasn't needed.

```go
type AccessTokenEntry struct {
	Resource     string    `json:"resource"`
	AccessToken  string    `json:"accessToken"`
	RefreshToken string    `json:"refreshToken"`
	ExpiresOn    time.Time `json:"expiresOn"`
}

type AccessTokenEntries []AccessTokenEntry

func (entry *AccessTokenEntry) UnmarshalJSON(data []byte) (err error) {

	// unmarshal as map
	var v map[string]interface{}
	if err = json.Unmarshal(data, &v); err != nil {
		return
	}

	// everything but expiresOn
	entry.Resource, _ = v["resource"].(string)
	entry.AccessToken, _ = v["accessToken"].(string)
	entry.RefreshToken, _ = v["refreshToken"].(string)

	// expiresOn
	expiresOnString, _ := v["expiresOn"].(string)
	expiresOnString = strings.Trim(expiresOnString, " ")
	if len(expiresOnString) > 0 {

		// parse the time
		var expiresOnParsed time.Time
		expiresOnParsed, err = time.Parse("2006-01-02 15:04:05.000000", expiresOnString)
		if err != nil {
			return
		}

		// the time is in local time even though it doesn't have a timezone, so fix it
		_, offset := time.Now().Zone()
		entry.ExpiresOn = expiresOnParsed.Local().Add(-1 * time.Second * time.Duration(offset))

	}

	return
}

func GetAccessToken(resource string) (accessToken string, err error) {

	// get the current user info
	usr, err := user.Current()
	if err != nil {
		return
	}

	// read the accessTokens.json file from disk
	content, err := ioutil.ReadFile(usr.HomeDir + "/.azure/accessTokens.json")
	if err != nil {
		return
	}

	// deserialize
	var entries AccessTokenEntries
	err = json.Unmarshal(content, &entries)
	if err != nil {
		return
	}

	// find the right resource
	for _, entry := range entries {
		if entry.Resource == resource {
			accessToken = entry.AccessToken
			if entry.ExpiresOn.Before(time.Now()) {
				// refresh
			}
			return
		}
	}

	// return not found
	err = fmt.Errorf("the specified resource (%s) was not found in the cache", resource)
	return

}
```

## Getting an Access Token by running az-cli

Pulling the cached access tokens from the file wasn't a good solution because I was going to have to check for expiry and fetch new ones, but using the az-cli to execute the command seemed a better solution...

```go
type AccessTokenEntry struct {
	AccessToken string `json:"accessToken"`
}

func GetAccessToken(ctx context.Context, resource string) (accessToken string, err error) {

	// execute the command and get the output
	cmd := exec.CommandContext(ctx, "az", "account", "get-access-token", "--resource", resource)
	var content []byte
	content, err = cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			err = errors.New(string(ee.Stderr))
		}
		return
	}

	// deserialize
	var entry AccessTokenEntry
	err = json.Unmarshal(content, &entry)
	if err != nil {
		return
	}

	// check for an accessToken
	// NOTE: this probably won't happen because cmd.Output() will probably throw a non-zero code
	if len(entry.AccessToken) < 1 {
		err = fmt.Errorf("the access token could not be obtained from az-cli - %s", content)
	}

	// return the accessToken
	accessToken = entry.AccessToken
	return

}
```

## Using reflection

When using go-envconfig, I need to print all variables to the console. This used reflection to print everything in a struct...

```go
func Print(config interface{}) {

	typeOf := reflect.TypeOf(config)
	valueOf := reflect.ValueOf(&config)
	elem := valueOf.Elem()

	fmt.Printf("Configuration:")
	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)
		value := elem.Field(i)
		if field.Tag.Get("config") == "mask" {
			fmt.Printf(`  %s = "%v"`, field.Name, "(masked)")
		} else {
			fmt.Printf(`  %s = "%v"`, field.Name, value)
		}
	}

}
```

## Inheritance

```go
package main

import "fmt"

// NOTE: there is no constructor concept so we have to use a method like Init() and remember to call it
// NOTE: not implementing lookup isn't caught by the compiler so the user of the method will cause a panic
// NOTE: there is no way I could ever find to call a method on derived from base (maybe reflection)
//   for example, if you defined a more specific version of a func() I cannot seem to call it

type IBase interface {
	Init()
	Print() string
	lookup(input string) (output string)
}

type Base struct {
	IBase
	base *Base
	name string
}

func (base *Base) SetBase() {
	base.base = base
}

/*
func (base *Base) lookup(input string) (output string) {
	output = input
	return
}
*/

func (base *Base) Print() {
	fmt.Printf("from base [%s]: %s\n", base.lookup("key"), base.name)
}

func (base *Base) Init() {
	base.name = "Peter"
}

type Derived struct {
	Base
	nickname string
}

func (derived *Derived) Print() {
	fmt.Printf("from derived [%s]: %s (%s)\n", derived.lookup("key"), derived.name, derived.nickname)
}

func (derived *Derived) Init() {
	derived.SetBase()
	derived.name = "Peter"
	derived.nickname = "Pete"
}

func (derived *Derived) PrintAsBase() {

	// NOTE: invalid type assertion: derived.(*Base) (non-interface type *Derived on left)
	// base := derived.(*Base)

	// NOTE: no incarnation of this works either
	// base := Base(*derived)

	// NOTE: the SetBase() method is about the only way to get this to work
	derived.base.Print()

}

func main() {
	var derived Derived
	derived.Init()
	derived.Print()
	derived.PrintAsBase()
}
```

the following output...

```text
from derived [key]: Peter (Pete)
from base [key]: Peter
```
