package clcgo

import (
	"errors"
	"fmt"
)

const ServerURL = "servers/%s/%s"

type ServerService struct {
	client *Client
}

type Server struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	LocationID string `json:"locationId"`
}

func (s *ServerService) Get(id string) (*Server, error) {

	if id == "" {
		return nil, errors.New("The server needs an ID attribute to generate a URL")
	}

	url := fmt.Sprintf(ServerURL, s.client.user.AccountAlias, id)
	req, err := s.client.newRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	server := &Server{}
	err = s.client.executeRequest(req, server)

	if err != nil {
		return nil, err
	}

	return server, err
}
