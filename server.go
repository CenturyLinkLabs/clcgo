package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ServerCreationURL  = APIRoot + "/servers/%s"
	ServerURL          = ServerCreationURL + "/%s"
	PublicIPAddressURL = ServerURL + "/publicIPAddresses"
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

func (s Server) URL(a string) (string, error) {
	if s.ID == "" && s.uuidURI == "" {
		return "", errors.New("An ID field is required to get a server")
	} else if s.uuidURI != "" {
		return APIDomain + s.uuidURI, nil
	}

	return fmt.Sprintf(ServerURL, a, s.ID), nil
}

func (s Server) RequestForSave(a string) (Request, error) {
	url := fmt.Sprintf(ServerCreationURL, a)
	p, err := s.parametersForSave()
	if err != nil {
		return Request{}, err
	}

	return Request{URL: url, Parameters: p}, nil
}

func (s Server) parametersForSave() (interface{}, error) {
	// TODO Freak out when the combo of AccountAlias and Name is too long! Which
	// is programatically defined and I don't have the rules.
	// TODO Well, actually.... CPU, MemoryGB, and Type are required, too! Those
	// are good candidates for defaults, though.
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

func (i PublicIPAddress) RequestForSave(a string) (Request, error) {
	if i.Server.ID == "" {
		return Request{}, errors.New("A Server with an ID is required to add a Public IP Address")
	}

	url := fmt.Sprintf(PublicIPAddressURL, a, i.Server.ID)
	return Request{URL: url, Parameters: i}, nil
}

func (i PublicIPAddress) StatusFromResponse(r []byte) (*Status, error) {
	l := Link{}
	err := json.Unmarshal(r, &l)
	if err != nil {
		return nil, err
	}

	return &Status{URI: l.HRef}, nil
}
