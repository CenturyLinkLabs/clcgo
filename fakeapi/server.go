package fakeapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

// CreateFakeServer instantiates a new httptest.Server hooked the the correct
// routes and returning the expected fixture data.
func CreateFakeServer() *httptest.Server {
	m := http.NewServeMux()
	s := httptest.NewServer(m)

	m.Handle("/v2/authentication/login", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := struct {
			Username string
			Password string
		}{}
		s, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(s, &p)

		if p.Username == "user" && p.Password == "pass" {
			fmt.Fprintf(w, AuthenticationSuccessfulResponse)
		} else {
			http.Error(w, "{}", 400)
		}
	}))

	m.Handle("/v2/servers/ACME/server1", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer 1234ABCDEF" {
			http.Error(w, "", 401)
		} else {
			fmt.Fprintf(w, ServerResponse)
		}
	}))

	return s
}
