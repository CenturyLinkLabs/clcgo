package clcgo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO look at implementing a testEntity like is done with testSavable in the
// savable_entities_test.

func TestSuccessfulGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(serverURL, "AA", id)
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)

	r.registerHandler("GET", url, func(token string, req request) (string, error) {
		assert.Equal(t, "token", token)
		return fmt.Sprintf(`{"name": "testname", "id": "%s"}`, id), nil
	})

	s := Server{ID: id}
	err := c.getEntity(r, &s)
	assert.NoError(t, err)
	assert.Equal(t, "testname", s.Name)
}

func TestErroredURLInGetEntity(t *testing.T) {
	r := newTestRequestor()
	s := Server{}
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	err := c.getEntity(&r, &s)

	_, e := s.url("abc123")
	assert.EqualError(t, err, e.Error())
}

func TestErroredInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	s := Server{ID: id}
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	err := c.getEntity(r, &s)
	url := fmt.Sprintf(serverURL, "AA", id)

	assert.EqualError(t, err, fmt.Sprintf("There is no handler for GET '%s'", url))
}

func TestBadJSONInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(serverURL, "AA", id)

	r.registerHandler("GET", url, func(token string, req request) (string, error) {
		return ``, nil
	})

	s := Server{ID: id}
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	err := c.getEntity(r, &s)

	assert.EqualError(t, err, "unexpected end of JSON input")
}
