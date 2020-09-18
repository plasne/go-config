<!-- markdownlint-disable no-hard-tabs -->

# go-config

stuff goes here

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

```go
func TestGoEnvConfig(t *testing.T) {
	ctx := context.Background()

	type Config struct {
		Int          int      `env:"INT"`
		IntDefault   int      `env:"INT,default=7"`
		Array        []string `env:"ARRAY"`
		ArrayDefault []string `env:"ARRAY,default=cat,dog"`
		Multi        string   `env:"HISTORY_DB,DB"`
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

	// NOTE: I believe this should be 7, not 0 because "dog" is not an integer
	t.Run("Int_multiple_env", func(t *testing.T) {
		// NOTE: the use case here is the user can set a DB connstring for everything, or
		//  can set more specific connstrings for different subsystems
		var config Config
		lookuper := envconfig.MapLookuper(map[string]string{
			"DB": "connstring",
		})
		if err := envconfig.ProcessWith(ctx, &config, lookuper); err != nil {
			t.Errorf("could not get configuration")
		}
		e := "connstring"
		if !reflect.DeepEqual(config.Multi, e) {
			t.Errorf(`string Failed: expected "%v", got "%v"`, e, config.Multi)
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
    --- FAIL: TestGoEnvConfigInt/Int_multiple_env (0.00s)
        config_test.go:49: could not get configuration
        config_test.go:53: string Failed: expected "connstring", got ""
    --- FAIL: TestGoEnvConfigInt/Array_default_from_empty (0.00s)
        config_test.go:64: could not get configuration
        config_test.go:68: array Failed: expected [cat dog], got []
```

+ no way to do lookup

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
