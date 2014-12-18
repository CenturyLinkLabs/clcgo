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
	url := fmt.Sprintf(ServerURL, "AA", id)
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}

	r.registerGetHandler(url, func(token string, req Request) (string, error) {
		assert.Equal(t, "token", token)
		return fmt.Sprintf(`{"name": "testname", "id": "%s"}`, id), nil
	})

	s := Server{ID: id}
	err := getEntity(r, c, &s)
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
	err := getEntity(r, c, &s)
	url := fmt.Sprintf(ServerURL, "AA", id)

	assert.EqualError(t, err, fmt.Sprintf("There is no handler for the URL '%s'", url))
}

func TestBadJSONInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	id := "abc123"
	url := fmt.Sprintf(ServerURL, "AA", id)

	r.registerGetHandler(url, func(token string, req Request) (string, error) {
		return ``, nil
	})

	s := Server{ID: id}
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	err := getEntity(r, c, &s)

	assert.EqualError(t, err, "unexpected end of JSON input")
}

func TestStatusURL(t *testing.T) {
	s := Status{URI: "/v2/status/1234"}
	url, err := s.URL("AA")
	assert.NoError(t, err)
	assert.Equal(t, APIDomain+"/v2/status/1234", url)
}

func TestStatusHasSucceeded(t *testing.T) {
	s := Status{}
	assert.False(t, s.HasSucceeded())

	s.Status = "executing"
	assert.False(t, s.HasSucceeded())

	s.Status = "succeeded"
	assert.True(t, s.HasSucceeded())
}
