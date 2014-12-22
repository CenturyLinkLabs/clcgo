package clcgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSavable struct {
	CallbackForRequest func(string) (request, error)
	CallbackForStatus  func([]byte) (*Status, error)
}

type savableCreationParameters struct {
	Value string
}

func (s testSavable) requestForSave(a string) (request, error) {
	if s.CallbackForRequest != nil {
		return s.CallbackForRequest(a)
	}

	return request{
		URL:        "/server/creation/url",
		Parameters: savableCreationParameters{Value: "testSavable"},
	}, nil
}

func (s testSavable) statusFromResponse(r []byte) (*Status, error) {
	if s.CallbackForStatus != nil {
		return s.CallbackForStatus(r)
	}

	return &Status{URI: "example.com/savable_status"}, nil
}

func TestSuccessfulSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	p := savableCreationParameters{Value: "testSavable"}
	st := &Status{}
	s := testSavable{
		CallbackForRequest: func(a string) (request, error) {
			assert.Equal(t, "AA", a)
			return request{URL: "/servers", Parameters: p}, nil
		},
		CallbackForStatus: func(r []byte) (*Status, error) {
			assert.Equal(t, []byte(serverCreationSuccessfulResponse), r)
			return st, nil
		},
	}

	r.registerHandler("POST", "/servers", func(token string, req request) (string, error) {
		assert.Equal(t, "token", token)
		assert.Equal(t, p, req.Parameters)

		return serverCreationSuccessfulResponse, nil
	})

	status, err := saveEntity(r, c, &s)
	assert.NoError(t, err)
	assert.Equal(t, st, status)
}

func TestErroredRequestSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{
		CallbackForRequest: func(a string) (request, error) {
			return request{}, errors.New("Test Request Error")
		},
	}

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Test Request Error")
}

func TestErroredPostJSONSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{}

	r.registerHandler("POST", "/server/creation/url", func(token string, req request) (string, error) {
		return "", errors.New("Error from PostJSON")
	})

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Error from PostJSON")
}

func TestErorredStatusSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{
		CallbackForStatus: func(r []byte) (*Status, error) {
			return nil, errors.New("Test Status Error")
		},
	}

	r.registerHandler("POST", "/server/creation/url", func(token string, req request) (string, error) {
		return "response", nil
	})

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Test Status Error")
}
