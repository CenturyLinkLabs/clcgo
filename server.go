package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	ServerURL         = APIRoot + "/servers/%s/%s"
	ServerCreationURL = APIRoot + "/servers/%s"
)

type Server struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	GroupID        string `json:"groupId"`
	SourceServerID string `json:"sourceServerId"` // nonexistant in get, extract to creation params?
	CPU            int    `json:"cpu"`
	MemoryGB       int    `json:"memoryGB"` // is memoryMB in get, extract to creation params?
	Type           string `json:"type"`
}

type serverCreationResponse struct {
	Links []Link `json:"links"`
}

func (s Server) URL(a string) (string, error) {
	if s.ID == "" {
		return "", errors.New("An ID field is required to get a server")
	}

	return fmt.Sprintf(ServerURL, a, s.ID), nil
}

func (s Server) URLForSave(a string) (string, error) {
	return fmt.Sprintf(ServerCreationURL, a), nil
}

func (s Server) ParametersForSave() (interface{}, error) {
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

	for _, l := range scr.Links {
		if l.Rel == "status" {
			return &Status{URI: l.HRef}, nil
		}
	}

	return nil, errors.New("The creation response has no status link")
}
