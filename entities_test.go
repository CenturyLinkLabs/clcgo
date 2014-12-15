package clcgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func (r testRequestor) GetJSON(t string, url string) ([]byte, error) {
	callback, found := r.GetHandlers[url]
	if found {
		s, err := callback(t, url)
		return []byte(s), err
	}

	return nil, fmt.Errorf("There is no handler for the URL '%s'", url)
}

func TestSuccessfulGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(ServerURL, "AA", id)
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}

	r.registerGetHandler(url, func(token string, url string) (string, error) {
		assert.Equal(t, "token", token)
		return fmt.Sprintf(`{"name": "testname", "id": "%s"}`, id), nil
	})

	s := Server{ID: id}
	err := getEntity(&r, c, &s)
	assert.NoError(t, err)
	assert.Equal(t, "testname", s.Name)
}

func TestErroredURLInGetEntity(t *testing.T) {
	r := newTestRequestor()
	s := Server{}
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	err := getEntity(&r, c, &s)

	_, e := s.URL("abc123")
	assert.EqualError(t, err, e.Error())
}

func TestErroredInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	s := Server{ID: id}
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	err := getEntity(&r, c, &s)
	url := fmt.Sprintf(ServerURL, "AA", id)

	assert.EqualError(t, err, fmt.Sprintf("There is no handler for the URL '%s'", url))
}

func TestBadJSONInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(ServerURL, "AA", id)

	r.registerGetHandler(url, func(token string, url string) (string, error) {
		return ``, nil
	})

	s := Server{ID: id}
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	err := getEntity(&r, c, &s)

	assert.EqualError(t, err, "unexpected end of JSON input")
}
