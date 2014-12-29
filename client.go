package clcgo

import "encoding/json"

type Client struct {
	Credentials Credentials
}

func ClientFromCredentials(c Credentials) *Client {
	return &Client{Credentials: c}
}

func (c *Client) GetEntity(e entity) error {
	return c.getEntity(&clcRequestor{}, e)
}

func (c *Client) getEntity(r requestor, e entity) error {
	url, err := e.url(c.Credentials.AccountAlias)
	if err != nil {
		return err
	}
	j, err := r.GetJSON(c.Credentials.BearerToken, request{URL: url})
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &e)
}

func (c *Client) SaveEntity(e savableEntity) (*Status, error) {
	return c.saveEntity(clcRequestor{}, e)
}

func (c *Client) saveEntity(r requestor, e savableEntity) (*Status, error) {
	req, err := e.requestForSave(c.Credentials.AccountAlias)
	if err != nil {
		return nil, err
	}
	resp, err := r.PostJSON(c.Credentials.BearerToken, req)
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
