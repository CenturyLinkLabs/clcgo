package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testParameters struct {
	TestKey string
}

func startServer(uri string, handler func(w http.ResponseWriter, r *http.Request)) (string, *net.TCPListener) {
	laddr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:31981")
	if nil != err {
		panic("Couldn't start server")
	}
	listener, err := net.ListenTCP("tcp", laddr)
	if nil != err {
		panic("Couldn't start server")
	}

	h := http.NewServeMux()

	h.HandleFunc(uri, func(w http.ResponseWriter, r *http.Request) {
		handler(w, r)
	})

	go http.Serve(listener, h)
	return "http://127.0.0.1:31981", listener
}

func TestSuccessfulPostJSON(t *testing.T) {
	root, l := startServer("/example", func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("accepts"))

		s, _ := ioutil.ReadAll(r.Body)
		var p testParameters
		err := json.Unmarshal(s, &p)
		assert.NoError(t, err)
		assert.Equal(t, testParameters{"Testing"}, p)

		fmt.Fprintf(w, "Response Text")
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(root+"/example", testParameters{"Testing"})
	assert.NoError(t, err)

	responseString := string(response)
	assert.Equal(t, "Response Text", responseString)
}

func TestUnhandledStatusOnPostJSON(t *testing.T) {
	root, l := startServer("/example", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(root+"/example", testParameters{"Testing"})
	assert.Contains(t, string(response), "Bad Request")

	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Got an unexpected status code")
		assert.Equal(t, 400, reqErr.StatusCode)
	}
}

func TestSuccessfulGetJSON(t *testing.T) {
	root, l := startServer("/example", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Response Text")

		assert.Equal(t, "Bearer token", r.Header.Get("Authorization"))
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "application/json", r.Header.Get("accepts"))
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", root+"/example")
	assert.NoError(t, err)
	assert.Equal(t, "Response Text", string(response))
}

func TestErrored401GetJson(t *testing.T) {
	root, l := startServer("/example", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.GetJSON("token", root+"/example")
	assert.Contains(t, string(response), "Bad Request")

	reqErr, ok := err.(RequestError)
	if assert.True(t, ok) {
		assert.EqualError(t, reqErr, "Got an unexpected status code")
		assert.Equal(t, 400, reqErr.StatusCode)
	}
}
