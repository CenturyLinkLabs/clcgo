package clcgo

import "testing"
import "fmt"

type handlerCallback func(string, authParameters) string

type testRequestor struct {
	Handlers map[string]handlerCallback
}

func newTestRequestor() testRequestor {
	return testRequestor{Handlers: make(map[string]handlerCallback)}
}

// TODO: And a count for verification
func (r *testRequestor) registerHandler(url string, callback handlerCallback) {
	r.Handlers[url] = callback
}

func (r *testRequestor) PostJSON(url string, v interface{}) ([]byte, error) {
	callback, found := r.Handlers[url]
	if found {
		// TODO error checking on coercion and make this more flexible
		response := callback(url, v.(authParameters))
		return []byte(response), nil
	} else {
		return nil, fmt.Errorf("There is no handler for the URL '%s'", url)
	}
}

// TODO: Handles bad response codes for explosions and bad passwords

func TestFetchCredentialsWithGoodCredentials(t *testing.T) {
	r := newTestRequestor()

	r.registerHandler(AuthenticationURL, func(url string, parameters authParameters) string {
		if parameters.Username != "username" {
			t.Errorf("Expected Username to be username, got '%s'", parameters.Username)
		}

		if parameters.Password != "password" {
			t.Errorf("Expected Password to be password, got '%s'", parameters.Password)
		}

		return `{"bearerToken":"expected token"}`
	})

	credentials, err := fetchCredentials(&r, "username", "password")
	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if credentials.BearerToken != "expected token" {
		t.Errorf("Expected a BearerToken and got '%s'", credentials.BearerToken)
	}
}
