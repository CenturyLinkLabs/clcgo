package clcgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchCredentialsWithGoodCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler(AuthenticationURL, func(url string, parameters authParameters) (string, error) {
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

	r.registerHandler(AuthenticationURL, func(url string, parameters authParameters) (string, error) {
		return "Bad Request", RequestError{"There was a problem with the request", 400}
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	assert.EqualError(t, err, "There was a problem with your credentials")
	assert.Equal(t, Credentials{}, credentials)
}
