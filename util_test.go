package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParameters struct {
	TestKey string
}

func TestSuccessfulUnauthenticatedPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Empty(t, r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("accepts"))

		s, _ := ioutil.ReadAll(r.Body)
		var p testParameters
		err := json.Unmarshal(s, &p)
		assert.NoError(t, err)
		assert.Equal(t, testParameters{"Testing"}, p)

		fmt.Fprintf(w, "Response Text")
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON("", ts.URL, testParameters{"Testing"})
	assert.NoError(t, err)

	responseString := string(response)
	assert.Equal(t, "Response Text", responseString)
}

func TestSuccessfulAuthenticatedPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		fmt.Fprintf(w, "Response Text")
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	_, err := r.PostJSON("token", ts.URL, testParameters{})
	assert.NoError(t, err)
}

func TestUnauthorizedPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", 401)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	_, err := r.PostJSON("token", ts.URL, testParameters{})
	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Your bearer token was rejected")
		assert.Equal(t, 401, reqErr.StatusCode)
	}
}

func TestUnhandledStatusOnPostJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON("", ts.URL, testParameters{"Testing"})
	assert.Contains(t, string(response), "Bad Request")

	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Got an unexpected status code")
		assert.Equal(t, 400, reqErr.StatusCode)
	}
}

func TestSuccessfulGetJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response Text")

		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("accepts"))
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", ts.URL)
	assert.NoError(t, err)
	assert.Equal(t, "Response Text", string(response))
}

func TestUnauthorizedGetJSON(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Unauthorized", 401)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	_, err := r.GetJSON("token", ts.URL)
	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Your bearer token was rejected")
		assert.Equal(t, 401, reqErr.StatusCode)
	}
}

func TestErrored400GetJson(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	}))
	defer ts.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", ts.URL)
	assert.Contains(t, string(response), "Bad Request")

	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Got an unexpected status code")
		assert.Equal(t, 400, reqErr.StatusCode)
	}
}
