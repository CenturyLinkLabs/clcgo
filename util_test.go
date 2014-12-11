package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"testing"
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
		fmt.Fprintf(w, "Response Text")

		if r.Method != "POST" {
			t.Errorf("Expected request method to be of type 'POST', got '%s'", r.Method)
		}

		if a := r.Header.Get("Content-Type"); a != "application/json" {
			t.Errorf("Expected request content type to be 'application/json', got '%s'", a)
		}

		if a := r.Header.Get("accepts"); a != "application/json" {
			t.Errorf("Expected request content type to be 'application/json', got '%s'", a)
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
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(root+"/example", testParameters{"Testing"})
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	responseString := string(response)
	if responseString != "Response Text" {
		t.Errorf("Expected response to be 'Response Text', got '%s'", responseString)
	}
}

func TestUnhandledStatusOnPostJSON(t *testing.T) {
	root, l := startServer("/example", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})
	defer l.Close()

	r := &CLCRequestor{}
	response, err := r.PostJSON(root+"/example", testParameters{"Testing"})
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
