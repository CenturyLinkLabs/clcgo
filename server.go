package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	serverCreationURL      = apiRoot + "/servers/%s"
	serverURL              = serverCreationURL + "/%s"
	publicIPAddressURL     = serverURL + "/publicIPAddresses"
	serverActiveStatus     = "active"
	serverPausedPowerState = "paused"
)

// A Server can be used to either fetch an existing Server or provision and new
// one. To fetch, you must supply an ID value. For creation, there are numerous
// required values. The API documentation should be consulted.
type Server struct {
	uuidURI        string `json:"-"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	GroupID        string `json:"groupId"`
	Status         string `json:"status"`
	SourceServerID string `json:"sourceServerId"` // TODO: nonexistent in get, extract to creation params?
	CPU            int    `json:"cpu"`
	MemoryGB       int    `json:"memoryGB"` // TODO: memoryMB in get, extract to creation params?
	Type           string `json:"type"`
	Details        struct {
		PowerState  string `json:"powerState"`
		IPAddresses []struct {
			Public   string `json:"public"`
			Internal string `json:"internal"`
		} `json:"ipAddresses"`
	} `json:"details"`
}

type serverCreationResponse struct {
	Links []Link `json:"links"`
}

// A PublicIPAddress can be created and associated with an existing,
// provisioned Server. You must supply the associated Server object.
//
// You can also supply an optional slice of Port objects that will make the
// specified ports accessible at the address.
type PublicIPAddress struct {
	Server Server
	Ports  []Port `json:"ports"`
}

// A Port object specifies a network port that should be made available on a
// PublicIPAddress. It can only be used in conjunction with the PublicIPAddress
// resource.
type Port struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

func (s Server) IsActive() bool {
	return s.Status == serverActiveStatus && !s.IsPaused()
}

func (s Server) IsPaused() bool {
	return s.Details.PowerState == serverPausedPowerState
}

func (s Server) URL(a string) (string, error) {
	if s.ID == "" && s.uuidURI == "" {
		return "", errors.New("An ID field is required to get a server")
	} else if s.uuidURI != "" {
		return apiDomain + s.uuidURI, nil
	}

	return fmt.Sprintf(serverURL, a, s.ID), nil
}

func (s Server) RequestForSave(a string) (request, error) {
	url := fmt.Sprintf(serverCreationURL, a)
	return request{URL: url, Parameters: s}, nil
}

func (s *Server) StatusFromResponse(r []byte) (*Status, error) {
	scr := serverCreationResponse{}
	err := json.Unmarshal(r, &scr)
	if err != nil {
		return nil, err
	}

	sl, err := typeFromLinks("status", scr.Links)
	if err != nil {
		return nil, errors.New("The creation response has no status link")
	}

	il, err := typeFromLinks("self", scr.Links)
	if err != nil {
		return nil, errors.New("The creation response has no self link")
	}

	s.uuidURI = il.HRef

	return &Status{URI: sl.HRef}, nil
}

func (i PublicIPAddress) RequestForSave(a string) (request, error) {
	if i.Server.ID == "" {
		return request{}, errors.New("A Server with an ID is required to add a Public IP Address")
	}

	url := fmt.Sprintf(publicIPAddressURL, a, i.Server.ID)
	return request{URL: url, Parameters: i}, nil
}

func (i PublicIPAddress) StatusFromResponse(r []byte) (*Status, error) {
	l := Link{}
	err := json.Unmarshal(r, &l)
	if err != nil {
		return nil, err
	}

	return &Status{URI: l.HRef}, nil
}
