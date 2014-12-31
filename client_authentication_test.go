package clcgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSuccessfulGetCredentials(t *testing.T) {
	r := newTestRequestor()
	c := Client{Requestor: r}

	r.registerHandler("POST", authenticationURL, func(token string, req request) (string, error) {
		c, ok := req.Parameters.(Credentials)
		assert.True(t, ok)
		assert.Empty(t, token)
		assert.Equal(t, "username", c.Username)
		assert.Equal(t, "password", c.Password)

		return `{"bearerToken":"expected token"}`, nil
	})

	err := c.GetCredentials("username", "password")
	assert.NoError(t, err)
	assert.Equal(t, "expected token", c.Credentials.BearerToken)
}

func TestErorredGetCredentials(t *testing.T) {
	r := newTestRequestor()
	c := Client{Requestor: r}

	r.registerHandler("POST", authenticationURL, func(token string, req request) (string, error) {
		err := RequestError{Message: "There was a problem with the request", StatusCode: 400}
		return "Bad Request", err
	})

	err := c.GetCredentials("username", "password")
	assert.EqualError(t, err, "There was a problem with your credentials")
	assert.Empty(t, c.Credentials.BearerToken)
}
