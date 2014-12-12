package clcgo

import "testing"

func TestFetchCredentialsWithGoodCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler(AuthenticationURL, func(url string, parameters authParameters) (string, error) {
		if parameters.Username != "username" {
			t.Errorf("Expected Username to be username, got '%s'", parameters.Username)
		}

		if parameters.Password != "password" {
			t.Errorf("Expected Password to be password, got '%s'", parameters.Password)
		}

		return `{"bearerToken":"expected token"}`, nil
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if credentials.BearerToken != "expected token" {
		t.Errorf("Expected a BearerToken and got '%s'", credentials.BearerToken)
	}
}

func TestFetchCredentialsWithBadCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler(AuthenticationURL, func(url string, parameters authParameters) (string, error) {
		return "Bad Request", RequestError{"There was a problem with the request", 400}
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	e := Credentials{}
	if credentials != e {
		t.Errorf("Expected empty Credentials, got '%s'", credentials)
	}

	if err.Error() != "There was a problem with your credentials" {
		t.Errorf("Expected specific error message, got '%s'", err)
	}
}
