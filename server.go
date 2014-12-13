package clcgo

import (
	"errors"
	"fmt"
)

const ServerURL = APIRoot + "/servers/%s/%s"

type Server struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s Server) URL(a string) (string, error) {
	if s.ID == "" {
		return "", errors.New("The server needs an ID attribute to generate a URL")
	}

	return fmt.Sprintf(ServerURL, a, s.ID), nil
}
