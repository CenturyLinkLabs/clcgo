package clcgo_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"github.com/CenturyLinkLabs/clcgo"
	"github.com/CenturyLinkLabs/clcgo/fakeapi"
)

var (
	exampleDefaultHTTPClient *http.Client
	fakeAPIServer            *httptest.Server
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
	fakeAPIServer = fakeapi.CreateFakeServer()

	// Replace the clcgo.DefaultHTTPClient with one that will rewrite the
	// requests to go to the test server instead of production.
	exampleDefaultHTTPClient = clcgo.DefaultHTTPClient
	u, _ := url.Parse(fakeAPIServer.URL)
	clcgo.DefaultHTTPClient = &http.Client{Transport: exampleDomainRewriter{RewriteURL: u}}
}

func teardownExample() {
	fakeAPIServer.Close()
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
	// Server Name: Test Name
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
