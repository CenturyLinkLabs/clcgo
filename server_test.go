package clcgo

import (
	"testing"
)

func TestWorkingServerURL(t *testing.T) {
	s := Server{ID: "abc123"}
	u, err := s.URL()
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}
	if e := ApiRoot + "/servers/abc123"; u != e {
		t.Errorf("Expected URL to be '%s', was '%s'", e, u)
	}
}

func TestErroredServerURL(t *testing.T) {
	s := Server{}
	u, err := s.URL()
	if err == nil {
		t.Errorf("Expected an error, got nothing")
	} else {
		if e := "The server needs an ID attribute to generate a URL"; err.Error() != e {
			t.Errorf("Expected the error '%s', got nothing", e)
		}
	}
	if u != "" {
		t.Errorf("Expected empty URL, got '%s'", u)
	}
}
