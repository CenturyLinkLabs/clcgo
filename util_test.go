package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testParameters struct {
	TestKey string
}

func TestSuccessfulPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected request method to be of type 'POST', got '%s'", r.Method)
		}

		if a := r.Header.Get("Content-Type"); a != "application/json" {
			t.Errorf("Expected request content type to be 'application/json', got '%s'", a)
		}

		if a := r.Header.Get("accepts"); a != "application/json" {
			t.Errorf("Expected request accepts to be 'application/json', got '%s'", a)
		}

		s, _ := ioutil.ReadAll(r.Body)
		var p testParameters
		if err := json.Unmarshal(s, &p); err != nil {
			t.Errorf("Expected no error, got '%s'", err)
		}

		e := testParameters{"Testing"}
		if p != e {
			t.Errorf("Expected '%s' and '%s' to match", e, p)
		}

		fmt.Fprintf(w, "Response Text")
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(ts.URL, testParameters{"Testing"})
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	responseString := string(response)
	if responseString != "Response Text" {
		t.Errorf("Expected response to be 'Response Text', got '%s'", responseString)
	}
}

func TestUnhandledStatusOnPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(ts.URL, testParameters{"Testing"})
	if a := strings.TrimSpace(string(response)); a != "Bad Request" {
		t.Errorf("Expected response 'Bad Request', got '%s'", a)
	}

	reqErr, ok := err.(RequestError)
	if !ok {
		t.Error("Expected a RequestError")
	}
	if reqErr.Err != "Got an unexpected status code" {
		t.Errorf("Expected a specific message, got '%s'", reqErr.Err)
	}
	if reqErr.StatusCode != 400 {
		t.Errorf("Expected a 400 status code got '%d'", reqErr.StatusCode)
	}
}

func TestSuccessfulGetJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response Text")

		a := r.Header.Get("Authorization")
		if e := "Bearer token"; a != e {
			t.Errorf("Expected Authorization header '%s', got '%s'", e, a)
		}

		if r.Method != "GET" {
			t.Errorf("Expected request method to be of type 'GET', got '%s'", r.Method)
		}

		if a := r.Header.Get("Content-Type"); a != "application/json" {
			t.Errorf("Expected request content type to be 'application/json', got '%s'", a)
		}

		if a := r.Header.Get("accepts"); a != "application/json" {
			t.Errorf("Expected request accepts to be 'application/json', got '%s'", a)
		}
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", ts.URL)
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	responseString := string(response)
	if responseString != "Response Text" {
		t.Errorf("Expected response to be 'Response Text', got '%s'", responseString)
	}
}

func TestErrored401GetJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", ts.URL)

	if a := strings.TrimSpace(string(response)); a != "Bad Request" {
		t.Errorf("Expected response 'Bad Request', got '%s'", a)
	}

	reqErr, ok := err.(RequestError)
	if !ok {
		t.Error("Expected a RequestError")
	}
	if reqErr.Err != "Got an unexpected status code" {
		t.Errorf("Expected a specific message, got '%s'", reqErr.Err)
	}
	if reqErr.StatusCode != 400 {
		t.Errorf("Expected a 400 status code got '%d'", reqErr.StatusCode)
	}
}
