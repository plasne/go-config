package config

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/Azure/go-autorest/autorest"
	"github.com/joho/godotenv"
)

// Allows a single line pattern that would emulate (condition ? true : false).
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

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

func GetAccessTokenForManagement(ctx context.Context) (string, error) {
	return GetAccessToken(ctx, "https://management.core.windows.net/")
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

	// TODO: can context apply to autorest somehow?

	// get an access token
	// NOTE: this cannot end in a slash
	token, err := GetAccessToken(ctx, "https://pelasne-config.azconfig.io")
	if err != nil {
		return
	}

	// request each filter
	for _, filter := range filters {

		// setup the request
		client := &autorest.Client{}
		q := map[string]interface{}{"key": filter}
		h := map[string]interface{}{"Authorization": "Bearer " + token}
		var req *http.Request
		req, err = autorest.Prepare(&http.Request{},
			autorest.AsGet(),
			autorest.WithBaseURL("https://pelasne-config.azconfig.io"),
			autorest.WithPath("/kv"),
			autorest.WithQueryParameters(q),
			autorest.WithHeaders(h))
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

	// get an access token
	// NOTE: this cannot end in a slash
	token, err := GetAccessToken(ctx, "https://vault.azure.net")
	if err != nil {
		return
	}

	// setup the request
	client := &autorest.Client{}
	q := map[string]interface{}{"api-version": "7.0"}
	h := map[string]interface{}{"Authorization": "Bearer " + token}
	var req *http.Request
	req, err = autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(url),
		autorest.WithQueryParameters(q),
		autorest.WithHeaders(h))
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

func main() {
	ctx := context.Background()
	// TODO: learn more about context

	// load the config
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}

	err = Apply(ctx, []string{"override:*", "sample:*"})
	if err != nil {
		panic(err)
	}

	AsString().TrySetByEnv("SECRET").Resolve(ctx).Print()

}
