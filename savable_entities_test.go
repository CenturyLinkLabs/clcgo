package clcgo

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testSavable struct {
	CallbackForURL        func(string) (string, error)
	CallbackForParameters func() (interface{}, error)
	CallbackForStatus     func([]byte) (*Status, error)
}

type savableCreationParameters struct {
	Value string
}

func (s testSavable) URLForSave(a string) (string, error) {
	if s.CallbackForURL != nil {
		return s.CallbackForURL(a)
	}

	return "/server/creation/url", nil
}

func (s testSavable) ParametersForSave() (interface{}, error) {
	if s.CallbackForParameters != nil {
		return s.CallbackForParameters()
	}

	return savableCreationParameters{Value: "testSavable"}, nil
}

func (s testSavable) StatusFromResponse(r []byte) (*Status, error) {
	if s.CallbackForStatus != nil {
		return s.CallbackForStatus(r)
	}

	return &Status{URL: "example.com/savable_status"}, nil
}

func TestSuccessfulSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	p := savableCreationParameters{Value: "testSavable"}
	st := &Status{}
	s := testSavable{
		CallbackForURL: func(a string) (string, error) {
			assert.Equal(t, "AA", a)
			return "/servers", nil
		},
		CallbackForParameters: func() (interface{}, error) {
			return p, nil
		},
		CallbackForStatus: func(r []byte) (*Status, error) {
			assert.Equal(t, []byte(serverCreationSuccessfulResponse), r)
			return st, nil
		},
	}

	r.registerHandler("/servers", func(token string, url string, v interface{}) (string, error) {
		assert.Equal(t, "token", token)
		assert.Equal(t, p, v)

		return serverCreationSuccessfulResponse, nil
	})

	status, err := saveEntity(r, c, &s)
	assert.NoError(t, err)
	assert.Equal(t, st, status)
}

func TestErroredURLSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{
		CallbackForURL: func(a string) (string, error) {
			return "", errors.New("Test URL Error")
		},
	}

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Test URL Error")
}

func TestErroredParametersSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{
		CallbackForParameters: func() (interface{}, error) {
			return nil, errors.New("Test Parameters Error")
		},
	}

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Test Parameters Error")
}

func TestErroredPostJSONSaveEntity(t *testing.T) {
	r := newTestRequestor()
	c := Credentials{BearerToken: "token", AccountAlias: "AA"}
	s := testSavable{}

	r.registerHandler("/server/creation/url", func(token string, url string, v interface{}) (string, error) {
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

	r.registerHandler("/server/creation/url", func(token string, url string, v interface{}) (string, error) {
		return "response", nil
	})

	status, err := saveEntity(r, c, &s)
	assert.Nil(t, status)
	assert.EqualError(t, err, "Test Status Error")
}
