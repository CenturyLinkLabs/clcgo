package clcgo

import "fmt"

type getHandlerCallback func(string, Request) (string, error)
type handlerCallback func(string, Request) (string, error)

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

func (r testRequestor) GetJSON(t string, req Request) ([]byte, error) {
	callback, found := r.GetHandlers[req.URL]
	if found {
		s, err := callback(t, req)
		return []byte(s), err
	}

	return nil, fmt.Errorf("There is no handler for the URL '%s'", req.URL)
}

func (r testRequestor) PostJSON(t string, req Request) ([]byte, error) {
	callback, found := r.Handlers[req.URL]
	if found {
		// TODO error checking on coercion and make this more flexible
		response, err := callback(t, req)
		return []byte(response), err
	}

	return nil, fmt.Errorf("There is no handler for the URL '%s'", req.URL)
}
