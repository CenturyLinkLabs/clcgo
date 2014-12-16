package clcgo

import "fmt"

type getHandlerCallback func(string, string) (string, error)
type handlerCallback func(string, string, interface{}) (string, error)

type testRequestor struct {
	GetHandlers map[string]getHandlerCallback
	Handlers    map[string]handlerCallback
}

func newTestRequestor() testRequestor {
	return testRequestor{
		Handlers:    make(map[string]handlerCallback),
		GetHandlers: make(map[string]getHandlerCallback),
	}
}

// TODO: And a count for verification
func (r *testRequestor) registerHandler(url string, callback handlerCallback) {
	r.Handlers[url] = callback
}

func (r *testRequestor) registerGetHandler(url string, callback getHandlerCallback) {
	r.GetHandlers[url] = callback
}

func (r testRequestor) PostJSON(t string, url string, v interface{}) ([]byte, error) {
	callback, found := r.Handlers[url]
	if found {
		// TODO error checking on coercion and make this more flexible
		response, err := callback(t, url, v)
		return []byte(response), err
	}

	return nil, fmt.Errorf("There is no handler for the URL '%s'", url)
}
