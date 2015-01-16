package clcgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

var creationResponse = `{"testSerializedKey":"value"}`

type testSavable struct {
	CallbackForRequest func(string) (request, error)
	TestSerializedKey  string `json:"testSerializedKey"`
}

type testStatusProviding struct {
	CallbackForStatus func([]byte) (Status, error)
}

type savableCreationParameters struct {
	Value string
}

func (s testSavable) RequestForSave(a string) (request, error) {
	if s.CallbackForRequest != nil {
		return s.CallbackForRequest(a)
	}

	return request{
		URL:        "/creation/url",
		Parameters: savableCreationParameters{Value: "testSavable"},
	}, nil
}

func (s testStatusProviding) RequestForSave(a string) (request, error) {
	return request{
		URL:        "/creation/url",
		Parameters: savableCreationParameters{Value: "testSavable"},
	}, nil
}

func (s testStatusProviding) StatusFromResponse(r []byte) (Status, error) {
	return s.CallbackForStatus(r)
}

func TestSuccessfulSavableSaveEntity(t *testing.T) {
	r := newTestRequestor()
	cr := APICredentials{BearerToken: "token", AccountAlias: "AA"}
	c := Client{Requestor: r, APICredentials: cr}
	p := savableCreationParameters{Value: "testSavable"}
	s := testSavable{
		CallbackForRequest: func(a string) (request, error) {
			assert.Equal(t, "AA", a)
			return request{URL: "/savable", Parameters: p}, nil
		},
	}

	r.registerHandler("POST", "/savable", func(token string, req request) (string, error) {
		assert.Equal(t, "token", token)
		assert.Equal(t, p, req.Parameters)

		return creationResponse, nil
	})

	status, err := c.SaveEntity(&s)
	assert.NoError(t, err)
	assert.True(t, status.HasSucceeded())
	assert.Equal(t, "value", s.TestSerializedKey)
}

func TestSuccessfulStatusProvidingSaveEntity(t *testing.T) {
	r := newTestRequestor()
	cr := APICredentials{BearerToken: "token", AccountAlias: "AA"}
	c := Client{Requestor: r, APICredentials: cr}
	st := Status{}
	s := testStatusProviding{
		CallbackForStatus: func(r []byte) (Status, error) {
			assert.Equal(t, []byte(creationResponse), r)
			return st, nil
		},
	}

	r.registerHandler("POST", "/creation/url", func(token string, req request) (string, error) {
		return creationResponse, nil
	})

	status, err := c.SaveEntity(&s)
	assert.NoError(t, err)
	assert.Equal(t, st, status)
}

func TestErroredRequestSaveEntity(t *testing.T) {
	r := newTestRequestor()
	cr := APICredentials{BearerToken: "token", AccountAlias: "AA"}
	c := Client{Requestor: r, APICredentials: cr}
	s := testSavable{
		CallbackForRequest: func(a string) (request, error) {
			return request{}, errors.New("Test Request Error")
		},
	}

	status, err := c.SaveEntity(&s)
	assert.Equal(t, Status{}, status)
	assert.EqualError(t, err, "Test Request Error")
}

func TestErroredPostJSONSaveEntity(t *testing.T) {
	r := newTestRequestor()
	cr := APICredentials{BearerToken: "token", AccountAlias: "AA"}
	c := Client{Requestor: r, APICredentials: cr}
	s := testSavable{}

	r.registerHandler("POST", "/creation/url", func(token string, req request) (string, error) {
		return "", errors.New("Error from PostJSON")
	})

	status, err := c.SaveEntity(&s)
	assert.Equal(t, Status{}, status)
	assert.EqualError(t, err, "Error from PostJSON")
}

func TestErorredStatusSaveEntity(t *testing.T) {
	r := newTestRequestor()
	cr := APICredentials{BearerToken: "token", AccountAlias: "AA"}
	c := Client{Requestor: r, APICredentials: cr}
	s := testStatusProviding{
		CallbackForStatus: func(r []byte) (Status, error) {
			return Status{}, errors.New("Test Status Error")
		},
	}

	r.registerHandler("POST", "/creation/url", func(token string, req request) (string, error) {
		return "response", nil
	})

	status, err := c.SaveEntity(&s)
	assert.Equal(t, Status{}, status)
	assert.EqualError(t, err, "Test Status Error")
}
