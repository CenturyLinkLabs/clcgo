package clcgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCredentialsWithGoodCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler("POST", authenticationURL, func(token string, req request) (string, error) {
		parameters, ok := req.Parameters.(authParameters)
		assert.True(t, ok)
		assert.Empty(t, token)
		assert.Equal(t, "username", parameters.Username)
		assert.Equal(t, "password", parameters.Password)

		return `{"bearerToken":"expected token"}`, nil
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	assert.NoError(t, err)
	assert.Equal(t, "expected token", credentials.BearerToken)
}

func TestFetchCredentialsWithBadCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler("POST", authenticationURL, func(token string, req request) (string, error) {
		err := RequestError{Message: "There was a problem with the request", StatusCode: 400}
		return "Bad Request", err
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	assert.EqualError(t, err, "There was a problem with your credentials")
	assert.Equal(t, Credentials{}, credentials)
}
