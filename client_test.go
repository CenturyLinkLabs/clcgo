package clcgo

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"
)

var (
	mux    *http.ServeMux
	server *httptest.Server
	client *Client
)

// setup sets up a test HTTP server and a hipchat.Client configured to talk
// to that test server.
// Tests should register handlers on mux which provide mock responses for
// the API method being tested.
func setup() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// github client configured to use test server
	client = NewClient("user", "pass")
	url, _ := url.Parse(server.URL)
	client.baseURL = url
}

// teardown closes the test HTTP server.
func teardown() {
	server.Close()
}

func TestNewRequestWithoutToken(t *testing.T) {
	c := NewClient("user", "pass")

	inURL, outURL := "foo", defaultBaseURL+"foo"
	inBody := struct {
		Message string `json:"message"`
	}{Message: "Hello"}
	outBody := `{"message":"Hello"}` + "\n"

	r, _ := c.newRequest("GET", inURL, inBody)

	if r.URL.String() != outURL {
		t.Errorf("NewRequest URL %s, want %s", r.URL.String(), outURL)
	}

	body, _ := ioutil.ReadAll(r.Body)
	if string(body) != outBody {
		t.Errorf("NewRequest body %s, want %s", body, outBody)
	}

	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("newRequest Content-Type header %s, want application/json", contentType)
	}

	accepts := r.Header.Get("Accepts")
	if accepts != "application/json" {
		t.Errorf("newRequest Accepts header %s, want application/json", accepts)
	}
}

func TestNewRequestWithToken(t *testing.T) {
	c := NewClient("user", "pass")
	c.user.BearerToken = "token"

	r, _ := c.newRequest("GET", "foo", nil)

	authorization := r.Header.Get("Authorization")
	if authorization != "Bearer "+c.user.BearerToken {
		t.Errorf("newRequest authorization header %s, want %s", authorization, "Bearer "+c.user.BearerToken)
	}
}

func TestExecuteRequest(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		Bar int
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprintf(w, `{"Bar":1}`)
	})
	req, _ := client.newRequest("GET", "/", nil)
	body := &foo{}

	err := client.executeRequest(req, body)

	if err != nil {
		t.Fatal(err)
	}

	want := &foo{Bar: 1}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}
