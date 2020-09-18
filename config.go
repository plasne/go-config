package config

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"

	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
)

// Allows a single line pattern that would emulate (condition ? true : false).
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

type authMode int

const (
	AuthMode_Env authMode = iota
	AuthMode_Cli
)

type preconfig struct {
	AUTH_MODE   authMode
	APPCONFIG   string
	CONFIG_KEYS []string
}

var config preconfig

func applyAuthorizer(client *autorest.Client, resource string) (err error) {

	// select
	var authorizer autorest.Authorizer
	switch config.AUTH_MODE {
	case AuthMode_Env:
		authorizer, err = auth.NewAuthorizerFromEnvironmentWithResource(resource)
	case AuthMode_Cli:
		authorizer, err = auth.NewAuthorizerFromCLIWithResource(resource)
	default:
		err = fmt.Errorf("there is no authorizer specified.")
		return
	}

	// check for errors
	if err != nil {
		return
	}

	// assign
	client.Authorizer = authorizer

	return
}

func tryExtractUrlForKeyvaultFromAppConfigEntry(value string) string {

	// make sure this is a keyvault entry
	lower := strings.ToLower(value)
	if !strings.HasPrefix(lower, "{") || !strings.HasSuffix(lower, "}") || !strings.Contains(lower, ".vault.azure.net/") {
		return value
	}

	// define the json pattern
	result := struct {
		Uri string `json:"uri"`
	}{}

	// deserialize
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		// ignore
		return value
	}

	// set the value to the Uri
	return result.Uri

}

func load(ctx context.Context, filters []string, useFullyQualifiedName bool) (values map[string]string, err error) {
	values = make(map[string]string)

	// make sure there is something to load
	if len(filters) < 1 {
		return
	}

	// make sure APPCONFIG is supplied so the load can happen
	if len(config.APPCONFIG) < 1 {
		err = fmt.Errorf("APPCONFIG was REQUIRED but not set.")
		return
	}

	// TODO: can context apply to autorest somehow?

	// request each filter
	for _, filter := range filters {

		// create/authorize the client
		client := &autorest.Client{}
		err = applyAuthorizer(client, config.APPCONFIG)
		if err != nil {
			return
		}

		// setup the request
		q := map[string]interface{}{"key": filter}
		var req *http.Request
		req, err = autorest.Prepare(&http.Request{},
			autorest.AsGet(),
			autorest.WithBaseURL(config.APPCONFIG),
			autorest.WithPath("/kv"),
			autorest.WithQueryParameters(q))
		if err != nil {
			return
		}

		// send the request
		var resp *http.Response
		resp, err = autorest.SendWithSender(client, req)
		if err != nil {
			return
		}

		// ensure it is something in the HTTP 200 range
		if resp.StatusCode < 200 || resp.StatusCode > 299 {
			err = fmt.Errorf("GET from appconfig (filter: %s) resulted in HTTP %d - %s", filter, resp.StatusCode, resp.Status)
			return
		}

		// define the json structure of the appconfig response
		result := struct {
			Items []struct {
				ContentType string `json:"content_type"`
				Key         string `json:"key"`
				Value       string `json:"value"`
			} `json:"items"`
		}{}

		// deserialize to json
		defer resp.Body.Close()
		dec := json.NewDecoder(resp.Body)
		err = dec.Decode(&result)
		if err != nil {
			return
		}

		// set the values
		for _, item := range result.Items {
			key := item.Key
			if !useFullyQualifiedName {
				path := strings.Split(item.Key, ":")
				key = strings.ToUpper(path[len(path)-1])
			}
			if _, ok := values[key]; !ok {
				val := tryExtractUrlForKeyvaultFromAppConfigEntry(item.Value)
				values[key] = val
			}
		}

	}

	return
}

func Load(ctx context.Context, filters []string) (values map[string]string, err error) {
	return load(ctx, filters, false)
}

func LoadFullyQualified(ctx context.Context, filters []string) (values map[string]string, err error) {
	return load(ctx, filters, true)
}

func Apply(ctx context.Context, filters []string) (err error) {

	// make sure there is something to apply
	if len(filters) < 1 {
		return
	}

	// load the values
	values, err := load(ctx, filters, false)
	if err != nil {
		return
	}

	// apply to env (if not already set)
	for key, value := range values {
		if _, ok := os.LookupEnv(key); !ok {
			os.Setenv(key, value)
		}
	}

	return
}

func resolve(ctx context.Context, url string) (val string, err error) {
	val = url

	// make sure this is a valid URL
	url = strings.ToLower(url)
	if !strings.HasPrefix(url, "https://") || !strings.Contains(url, ".vault.azure.net/") {
		return
	}

	// TODO: can context apply to autorest somehow?

	// create/authorize the client
	client := &autorest.Client{}
	err = applyAuthorizer(client, "https://vault.azure.net")
	if err != nil {
		return
	}

	// setup the request
	q := map[string]interface{}{"api-version": "7.0"}
	var req *http.Request
	req, err = autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(url),
		autorest.WithQueryParameters(q))
	if err != nil {
		return
	}

	// send the request
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req)
	if err != nil {
		return
	}

	// ensure it is something in the HTTP 200 range
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		err = fmt.Errorf("GET from keyvault (url: %s) resulted in HTTP %d - %s", url, resp.StatusCode, resp.Status)
		return
	}

	// define the json structure of the keyvault response
	result := struct {
		Value string `json:"value"`
	}{}

	// deserialize to json
	defer resp.Body.Close()
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return
	}

	// extract the value
	val = result.Value

	return
}

func ResolveAll(ctx context.Context, list []*StringChain) *sync.WaitGroup {
	wg := new(sync.WaitGroup)
	wg.Add(len(list))

	// resolve everything in parallel
	for _, chain := range list {
		go func(c *StringChain) {
			defer wg.Done()
			c.Resolve(ctx)
		}(chain)
	}

	return wg
}

func Startup(ctx context.Context) (err error) {

	// load from dotenv
	//  NOTE: ignore *os.PathError (the file is optional)
	err = godotenv.Load()
	if err != nil {
		if _, ok := err.(*os.PathError); !ok {
			return err
		} else {
			err = nil
		}
	}

	// do pre-configuration
	fmt.Println("PRE-CONFIGURATION:")
	table := map[string]int{
		"env": int(AuthMode_Env),
		"cli": int(AuthMode_Cli),
	}
	config.AUTH_MODE = authMode(AsInt().TrySetByEnv("AUTH_MODE").Lookup(table).DefaultTo(0).PrintLookup(table).Value())
	config.APPCONFIG = AsString().TrySetByEnv("APPCONFIG").Transform(func(chain *StringChain) {
		if chain.IsSet() {
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
			chain.SetTo(val)
		}
	}).Print().Value()
	config.CONFIG_KEYS = AsSplice().TrySetByEnv("CONFIG_KEYS").Print().Value()

	// load from appconfig
	if len(config.APPCONFIG) > 0 && len(config.CONFIG_KEYS) > 0 {
		err = Apply(ctx, config.CONFIG_KEYS)
		if err != nil {
			return
		}
	}

	return
}
