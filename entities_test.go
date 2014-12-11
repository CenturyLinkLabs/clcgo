package clcgo

import (
	"fmt"
	"testing"
)

func (r *testRequestor) GetJSON(url string) ([]byte, error) {
	callback, found := r.GetHandlers[url]
	if found {
		s, err := callback(url)
		return []byte(s), err
	} else {
		return nil, fmt.Errorf("There is no handler for the URL '%s'", url)
	}
}

func TestSuccessfulGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(ServerURL, id)

	r.registerGetHandler(url, func(url string) (string, error) {
		return fmt.Sprintf(`{"name": "testname", "id": "%s"}`, id), nil
	})

	s := Server{ID: id}
	//TODO audit ALL the pointer receivers
	err := getEntity(&r, "beartoken", &s)

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if e := "testname"; s.Name != e {
		t.Errorf("Expected Name to be '%s', got '%s'", e, s.Name)
	}
}

func TestErroredURLInGetEntity(t *testing.T) {
	r := newTestRequestor()
	s := Server{}
	err := getEntity(&r, "beartoken", &s)

	_, e := s.URL()
	if err.Error() != e.Error() {
		t.Errorf("Expected the error '%s', got '%s'", e, err)
	}
}

func TestErroredInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	s := Server{ID: id}
	err := getEntity(&r, "beartoken", &s)
	url := fmt.Sprintf(ServerURL, id)

	if e := fmt.Sprintf("There is no handler for the URL '%s'", url); err.Error() != e {
		t.Errorf("Expected the error '%s', got '%s'", e, err)
	}
}

func TestBadJSONInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(ServerURL, id)

	r.registerGetHandler(url, func(url string) (string, error) {
		return ``, nil
	})

	s := Server{ID: id}
	err := getEntity(&r, "beartoken", &s)

	if e := "unexpected end of JSON input"; err.Error() != e {
		t.Errorf("Expected the error '%s', got '%s'", e, err)
	}
}
