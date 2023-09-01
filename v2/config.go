package config

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
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

type preconfig struct {
	GOCONFIG_CREDS          []string
	GOCONFIG_APPCONFIG      string
	GOCONFIG_APPCONFIG_KEYS []string
}

var config preconfig
var credentialLock sync.Mutex
var credential *azidentity.ChainedTokenCredential
var tokenLock sync.Mutex
var tokens map[string]azcore.AccessToken = make(map[string]azcore.AccessToken)
var sharedHttpTransport *http.Transport

func createSharedHttpTransport() *http.Transport {
	defaultTransport := http.DefaultTransport.(*http.Transport)
	return &http.Transport{
		Proxy:                 defaultTransport.Proxy,
		DialContext:           defaultTransport.DialContext,
		MaxIdleConns:          defaultTransport.MaxIdleConns,
		IdleConnTimeout:       defaultTransport.IdleConnTimeout,
		TLSHandshakeTimeout:   defaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout: defaultTransport.ExpectContinueTimeout,
		TLSClientConfig: &tls.Config{
			MinVersion:    tls.VersionTLS12,
			Renegotiation: tls.RenegotiateNever,
		},
	}
}

func GetCredential() (azcore.TokenCredential, error) {
	credentialLock.Lock()
	defer credentialLock.Unlock()

	// check cache
	if credential != nil {
		return credential, nil
	}

	// create a list of creds
	var creds []azcore.TokenCredential
	for _, cred := range config.GOCONFIG_CREDS {
		switch strings.ToLower(cred) {
		case "env":
			env, err := azidentity.NewEnvironmentCredential(nil)
			if err == nil {
				creds = append(creds, env)
			}
		case "mi":
			mi, err := azidentity.NewManagedIdentityCredential(nil)
			if err == nil {
				creds = append(creds, mi)
			}
		case "cli":
			cli, err := azidentity.NewAzureCLICredential(nil)
			if err == nil {
				creds = append(creds, cli)
			}
		case "default":
			dft, err := azidentity.NewDefaultAzureCredential(nil)
			if err == nil {
				creds = append(creds, dft)
			}
		default:
			return nil, fmt.Errorf("GOCONFIG_CREDS contained an unsupported value: %s", cred)
		}
	}

	// create the chain
	// TODO: reimplement chain with mi having a timeout
	chain, err := azidentity.NewChainedTokenCredential(creds, nil)
	if err != nil {
		return nil, err
	}

	credential = chain
	return chain, nil
}

func GetAccessToken(ctx context.Context, scope string) (string, error) {
	tokenLock.Lock()
	defer tokenLock.Unlock()

	// check cache
	token, ok := tokens[scope]
	if ok && time.Until(token.ExpiresOn).Minutes() >= 5 {
		return token.Token, nil
	}

	// get credential
	cred, err := GetCredential()
	if err != nil {
		return "", err
	}

	// get token
	opt := policy.TokenRequestOptions{Scopes: []string{scope}}
	accessToken, err := cred.GetToken(ctx, opt)
	if err != nil {
		return "", err
	}

	tokens[scope] = accessToken
	return accessToken.Token, nil
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
	if len(config.GOCONFIG_APPCONFIG) < 1 {
		err = fmt.Errorf("GOCONFIG_APPCONFIG was REQUIRED but not set")
		return
	}

	// request each filter
	// TODO: improve performance by fetching these concurrently
	for _, filter := range filters {
		// create the client
		client := &autorest.Client{
			Sender: &http.Client{Transport: sharedHttpTransport},
		}

		// get the token
		var token string
		token, err = GetAccessToken(ctx, "https://azconfig.io")
		if err != nil {
			return
		}

		// setup the request
		q := map[string]interface{}{"key": filter}
		var req *http.Request
		req, err = autorest.Prepare(&http.Request{},
			autorest.AsGet(),
			autorest.WithBaseURL(config.GOCONFIG_APPCONFIG),
			autorest.WithPath("/kv"),
			autorest.WithQueryParameters(q),
			autorest.WithBearerAuthorization(token))
		if err != nil {
			return
		}

		// send the request
		var resp *http.Response
		resp, err = autorest.SendWithSender(client, req.WithContext(ctx))
		if err != nil {
			return
		}
		defer resp.Body.Close()

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

	// create the client
	client := &autorest.Client{
		Sender: &http.Client{Transport: sharedHttpTransport},
	}

	// get the token
	var token string
	token, err = GetAccessToken(ctx, "https://vault.azure.net")
	if err != nil {
		return
	}

	// setup the request
	q := map[string]interface{}{"api-version": "7.4"}
	var req *http.Request
	req, err = autorest.Prepare(&http.Request{},
		autorest.AsGet(),
		autorest.WithBaseURL(url),
		autorest.WithQueryParameters(q),
		autorest.WithBearerAuthorization(token))
	if err != nil {
		return
	}

	// send the request
	var resp *http.Response
	resp, err = autorest.SendWithSender(client, req.WithContext(ctx))
	if err != nil {
		return
	}
	defer resp.Body.Close()

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
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&result)
	if err != nil {
		return
	}

	// extract the value
	val = result.Value

	return
}

func areSlicesEqual(a []string, b []string) bool {
	if len(a) == len(b) {

		// shortcut as they are both 0
		if len(a) == 0 {
			return true
		}

		// check the values
		for i := range a {
			if a[i] != b[i] {
				return false
			}
		}

		return true
	} else {
		return false
	}
}

func Startup(ctx context.Context) (err error) {

	// create a shared http transport
	if sharedHttpTransport == nil {
		sharedHttpTransport = createSharedHttpTransport()
	}

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
	config.GOCONFIG_CREDS = AsSlice().TrySetByEnv("GOCONFIG_CREDS").DefaultTo([]string{"default"}).Print().Value()
	config.GOCONFIG_APPCONFIG = AsString().TrySetByEnv("GOCONFIG_APPCONFIG").Transform(func(chain *StringChain) {
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
	config.GOCONFIG_APPCONFIG_KEYS = AsSlice().TrySetByEnv("GOCONFIG_APPCONFIG_KEYS").Print().Value()

	// load from appconfig
	if len(config.GOCONFIG_APPCONFIG) > 0 && len(config.GOCONFIG_APPCONFIG_KEYS) > 0 {
		err = Apply(ctx, config.GOCONFIG_APPCONFIG_KEYS)
		if err != nil {
			return
		}
	}

	return
}
