package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const ServerURL = ApiRoot + "/servers/%s/%s"

type Server struct {
	ID string
	// TODO Deal with AccountAlias and credentials
	AccountAlias string
	Name         string
}

func (s Server) URL() (string, error) {
	if s.ID == "" || s.AccountAlias == "" {
		return "", errors.New("The server needs an AccountAlias and ID attribute to generate a URL")
	}

	return fmt.Sprintf(ServerURL, s.AccountAlias, s.ID), nil
}

func (s *Server) Unmarshal(j []byte) error {
	return json.Unmarshal(j, s)
}
