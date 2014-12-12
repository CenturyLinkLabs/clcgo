package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

const ServerURL = ApiRoot + "/servers/%s/%s"

type Server struct {
	ID   string
	Name string
}

func (s Server) URL(a string) (string, error) {
	if s.ID == "" {
		return "", errors.New("The server needs an ID attribute to generate a URL")
	}

	return fmt.Sprintf(ServerURL, a, s.ID), nil
}

func (s *Server) Unmarshal(j []byte) error {
	return json.Unmarshal(j, s)
}
