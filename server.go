package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const ServerURL = ApiRoot + "/servers/%s/%s"

type Server struct {
	Credentials Credentials
	ID          string
	Name        string
}

func (s *Server) SetCredentials(c Credentials) {
	s.Credentials = c
}

func (s Server) URL() (string, error) {
	if s.ID == "" {
		return "", errors.New("The server needs an ID attribute to generate a URL")
	}

	return fmt.Sprintf(ServerURL, s.Credentials.AccountAlias, s.ID), nil
}

func (s *Server) Unmarshal(j []byte) error {
	return json.Unmarshal(j, s)
}
