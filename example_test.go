package clcgo_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/CenturyLinkLabs/clcgo"
)

var (
	exampleDefaultHTTPClient *http.Client
	exampleHTTPServer        *httptest.Server
)

type exampleDomainRewriter struct {
	RewriteURL *url.URL
}

func (r exampleDomainRewriter) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	req.URL.Scheme = "http"
	req.URL.Host = r.RewriteURL.Host

	t := http.Transport{}
	return t.RoundTrip(req)
}

func setupExample() {
	// Set up a fake API server.
	m := http.NewServeMux()
	exampleHTTPServer = httptest.NewServer(m)

	m.Handle("/v2/authentication/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := struct {
			Username string
			Password string
		}{}
		s, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(s, &p)

		if p.Username == "user" && p.Password == "pass" {
			json := `{ "bearerToken": "1234ABCDEF", "accountAlias": "ACME" }`
			fmt.Fprintf(w, json)
		} else {
			http.Error(w, "{}", 400)
		}
	}))

	m.Handle("/v2/servers/ACME/server1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer 1234ABCDEF" {
			http.Error(w, "", 401)
		} else {
			json := `{
				"id": "server1",
				"name": "Test Server",
				"groupId": "group1",
				"details": {
					"ipAddresses": [
					{
						"internal": "10.0.0.1"
					},
					{
						"public": "8.8.8.8",
						"internal": "10.0.0.2"
					}
					]
				}
			}`
			fmt.Fprintf(w, json)
		}
	}))

	// Replace the clcgo.DefaultHTTPClient with one that will rewrite the
	// requests to go to the test server instead of production.
	exampleDefaultHTTPClient = clcgo.DefaultHTTPClient
	u, _ := url.Parse(exampleHTTPServer.URL)
	clcgo.DefaultHTTPClient = &http.Client{Transport: exampleDomainRewriter{RewriteURL: u}}
}

func teardownExample() {
	exampleHTTPServer.Close()
	clcgo.DefaultHTTPClient = exampleDefaultHTTPClient
}

func ExampleClient_GetCredentials_successful() {
	// Some test-related setup code which you can safely ignore.
	setupExample()
	defer teardownExample()

	c := clcgo.NewClient()
	c.GetCredentials("user", "pass")

	fmt.Printf("Account Alias: %s", c.Credentials.AccountAlias)
	// Output:
	// Account Alias: ACME
}

func ExampleClient_GetCredentials_failed() {
	// Some test-related setup code which you can safely ignore.
	setupExample()
	defer teardownExample()

	c := clcgo.NewClient()
	err := c.GetCredentials("bad", "bad")

	fmt.Printf("Error: %s", err)
	// Output:
	// Error: There was a problem with your credentials
}

func ExampleClient_GetEntity_successful() {
	// Some test-related setup code which you can safely ignore.
	setupExample()
	defer teardownExample()

	c := clcgo.NewClient()
	c.GetCredentials("user", "pass")

	s := clcgo.Server{ID: "server1"}
	c.GetEntity(&s)

	fmt.Printf("Server Name: %s", s.Name)
	// Output:
	// Server Name: Test Server
}

func ExampleClient_GetEntity_expiredToken() {
	// Some test-related setup code which you can safely ignore.
	setupExample()
	defer teardownExample()

	c := clcgo.NewClient()
	// You are caching this Bearer Token value and it has either expired or for
	// some other reason become invalid.
	c.Credentials = clcgo.Credentials{BearerToken: "expired", AccountAlias: "ACME"}

	s := clcgo.Server{ID: "server1"}
	err := c.GetEntity(&s)

	rerr, _ := err.(clcgo.RequestError)
	fmt.Printf("Error: %s, Status Code: %d", rerr, rerr.StatusCode)
	// Output:
	// Error: Your bearer token was rejected, Status Code: 401
}
