package clcgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var entityResponse = `{"testSerializedKey":"value"}`

type testEntity struct {
	CallbackForURL    func(string) (string, error)
	TestSerializedKey string `json:"testSerializedKey"`
}

func (e testEntity) url(a string) (string, error) {
	if e.CallbackForURL != nil {
		return e.CallbackForURL(a)
	}

	return "/entity/url", nil
}

func TestSuccessfulGetEntity(t *testing.T) {
	r := newTestRequestor()
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)

	r.registerHandler("GET", "/entity", func(token string, req request) (string, error) {
		assert.Equal(t, "token", token)
		return entityResponse, nil
	})

	e := testEntity{
		CallbackForURL: func(a string) (string, error) {
			assert.Equal(t, "AA", a)
			return "/entity", nil
		},
	}

	err := c.getEntity(r, &e)
	assert.NoError(t, err)
	assert.Equal(t, "value", e.TestSerializedKey)
}

func TestErroredURLInGetEntity(t *testing.T) {
	r := newTestRequestor()
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	e := testEntity{
		CallbackForURL: func(a string) (string, error) {
			return "", errors.New("Test URL Error")
		},
	}

	err := c.getEntity(&r, &e)
	assert.EqualError(t, err, "Test URL Error")
}

func TestErroredInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	e := testEntity{}
	r.registerHandler("GET", "/entity/url", func(token string, req request) (string, error) {
		return "", errors.New("Error from GetJSON")
	})

	err := c.getEntity(r, &e)
	assert.EqualError(t, err, "Error from GetJSON")
}

func TestBadJSONInGetJSONInGetEntity(t *testing.T) {
	r := newTestRequestor()
	cr := Credentials{BearerToken: "token", AccountAlias: "AA"}
	c := ClientFromCredentials(cr)
	e := testEntity{}
	r.registerHandler("GET", "/entity/url", func(token string, req request) (string, error) {
		return ``, nil
	})

	err := c.getEntity(r, &e)
	assert.EqualError(t, err, "unexpected end of JSON input")
}
