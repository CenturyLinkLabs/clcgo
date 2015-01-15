package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

type PauseServer struct {
	Server Server
}

type operationResponse struct {
	Links []Link `json:"links"`
}

const (
	operationsRoot = apiRoot + "/operations"
	pauseURL       = operationsRoot + "/%s/servers/pause"
)

func (p PauseServer) RequestForSave(a string) (request, error) {
	if p.Server.ID == "" {
		return request{}, errors.New("PauseServer requires a Server with an ID to pause!")
	}

	r := request{
		URL:        fmt.Sprintf(pauseURL, a),
		Parameters: []string{p.Server.ID},
	}

	return r, nil
}

func (p PauseServer) StatusFromResponse(r []byte) (*Status, error) {
	ors := []operationResponse{}
	err := json.Unmarshal(r, &ors)
	if err != nil {
		return nil, err
	}

	// This API is capable of operating on multiple servers in one call, but
	// allowing for single or multiple entities everywhere is going to require
	// more reshuffling than we need to do right now. We enforce a single server
	// on the operation, and error if the API returns more than one operation. It
	// shouldn't based on us only submitting one server a time, but just in case.
	// I'd rather error than ignore or panic.
	if len(ors) != 1 {
		return nil, errors.New("Expected a single operation response from the API!")
	}
	or := ors[0]

	sl, err := typeFromLinks("status", or.Links)
	if err != nil {
		return nil, errors.New("The operation response has no status link")
	}

	return &Status{URI: sl.HRef}, nil
}
