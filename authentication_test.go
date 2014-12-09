package clcgo

import "testing"

type testRequestor struct {
	Calls []call
}

func newTestRequestor() testRequestor {
	return testRequestor{Calls: make([]call, 0)}
}

type call struct {
	URL        string
	Parameters authParameters
}

func (r *testRequestor) PostJSON(url string, v interface{}) ([]byte, error) {
	parameters := v.(authParameters)
	r.Calls = append(r.Calls, call{url, parameters})

	return []byte(`{"bearerToken":"expected token"}`), nil
}

// TODO: Handles bad response codes for explosions and bad passwords

func TestFetchCredentialsWithGoodCredentials(t *testing.T) {
	r := newTestRequestor()
	credentials, err := fetchCredentials(&r, "username", "password")
	if err != nil {
		t.Errorf("Expected no error, got %s", err.Error)
	}

	if len(r.Calls) != 1 {
		t.Errorf("Expected PostJSON called once, callse %d times", len(r.Calls))
	}

	c := r.Calls[0]

	if c.URL != AuthenticationURL {
		t.Errorf("Expected call to '%s', got '%s'", AuthenticationURL, c.URL)
	}

	if c.Parameters.Username != "username" {
		t.Errorf("Expected Username to be username, got '%s'", c.Parameters.Username)
	}

	if c.Parameters.Password != "password" {
		t.Errorf("Expected Password to be password, got '%s'", c.Parameters.Password)
	}

	if credentials.BearerToken != "expected token" {
		t.Errorf("Expected a BearerToken and got '%s'", credentials.BearerToken)
	}
}
