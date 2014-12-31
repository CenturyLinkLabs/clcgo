package clcgo

import (
	"encoding/json"
	"errors"
)

const authenticationURL = apiRoot + "/authentication/login"

type Credentials struct {
	BearerToken  string `json:"bearerToken"`
	AccountAlias string `json:"accountAlias"`
	Username     string `json:"username"` // TODO: nonexistant in get, extract to creation params?
	Password     string `json:"password"` // TODO: nonexistant in get, extract to creation params?
}

type Client struct {
	Credentials Credentials
	Requestor   requestor
}

func (c Credentials) requestForSave(a string) (request, error) {
	return request{URL: authenticationURL, Parameters: c}, nil
}

func NewClient() *Client {
	return &Client{Requestor: clcRequestor{}}
}

func (c *Client) GetCredentials(u string, p string) error {
	c.Credentials = Credentials{Username: u, Password: p}
	_, err := c.SaveEntity(&c.Credentials)
	if err != nil {
		if rerr, ok := err.(RequestError); ok && rerr.StatusCode == 400 {
			err = errors.New("There was a problem with your credentials")
		}

		return err
	}

	return nil
}

func (c *Client) GetEntity(e entity) error {
	url, err := e.url(c.Credentials.AccountAlias)
	if err != nil {
		return err
	}
	j, err := c.Requestor.GetJSON(c.Credentials.BearerToken, request{URL: url})
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &e)
}

func (c *Client) SaveEntity(e savableEntity) (*Status, error) {
	req, err := e.requestForSave(c.Credentials.AccountAlias)
	if err != nil {
		return nil, err
	}
	resp, err := c.Requestor.PostJSON(c.Credentials.BearerToken, req)
	if err != nil {
		return nil, err
	}

	if spe, ok := e.(statusProvidingEntity); ok {
		status, err := spe.statusFromResponse(resp)
		if err != nil {
			return nil, err
		}

		return status, nil
	}

	json.Unmarshal(resp, &e)
	return &Status{Status: successfulStatus}, nil
}
