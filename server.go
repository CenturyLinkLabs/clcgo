package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	serverCreationURL  = apiRoot + "/servers/%s"
	serverURL          = serverCreationURL + "/%s"
	publicIPAddressURL = serverURL + "/publicIPAddresses"
)

type Server struct {
	uuidURI        string `json:"-"`
	ID             string `json:"id"`
	Name           string `json:"name"`
	GroupID        string `json:"groupId"`
	SourceServerID string `json:"sourceServerId"` // TODO: nonexistant in get, extract to creation params?
	CPU            int    `json:"cpu"`
	MemoryGB       int    `json:"memoryGB"` // TODO: memoryMB in get, extract to creation params?
	Type           string `json:"type"`
	Details        struct {
		IPAddresses []struct {
			Public   string `json:"public"`
			Internal string `json:"internal"`
		} `json:"ipAddresses"`
	} `json:"details"`
}

type serverCreationResponse struct {
	Links []Link `json:"links"`
}

type PublicIPAddress struct {
	Server Server
	Ports  []Port `json:"ports"`
}

type Port struct {
	Protocol string `json:"protocol"`
	Port     int    `json:"port"`
}

func (s Server) url(a string) (string, error) {
	if s.ID == "" && s.uuidURI == "" {
		return "", errors.New("An ID field is required to get a server")
	} else if s.uuidURI != "" {
		return apiDomain + s.uuidURI, nil
	}

	return fmt.Sprintf(serverURL, a, s.ID), nil
}

func (s Server) requestForSave(a string) (request, error) {
	url := fmt.Sprintf(serverCreationURL, a)
	p, err := s.parametersForSave()
	if err != nil {
		return request{}, err
	}

	return request{URL: url, Parameters: p}, nil
}

// BUG(dp): A combination of long AccountAlias and Name can be invalid. We need
// to decide whether this is worth verifying here or not.
// BUG(dp): There are additional fields that are also required. We need to
// decide how much verification is worthwhile. Those fields are CPU, MemoryGB,
// and Type.
func (s Server) parametersForSave() (interface{}, error) {
	f := []string{s.GroupID, s.Name, s.SourceServerID}
	for _, v := range f {
		if v == "" {
			return nil, errors.New(
				"The following fields are required to save: Name, GroupID, SourceServerID",
			)
		}
	}

	return s, nil
}

func (s *Server) statusFromResponse(r []byte) (*Status, error) {
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

func (i PublicIPAddress) requestForSave(a string) (request, error) {
	if i.Server.ID == "" {
		return request{}, errors.New("A Server with an ID is required to add a Public IP Address")
	}

	url := fmt.Sprintf(publicIPAddressURL, a, i.Server.ID)
	return request{URL: url, Parameters: i}, nil
}

func (i PublicIPAddress) statusFromResponse(r []byte) (*Status, error) {
	l := Link{}
	err := json.Unmarshal(r, &l)
	if err != nil {
		return nil, err
	}

	return &Status{URI: l.HRef}, nil
}
