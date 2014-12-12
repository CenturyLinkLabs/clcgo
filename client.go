package clcgo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const (
	defaultBaseURL = "https://api.tier3.com/v2/"
)

type user struct {
	BearerToken  string `json:"bearerToken"`
	AccountAlias string `json:"accountAlias"`
}

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Client struct {
	credentials credentials
	user        user
	baseURL     *url.URL
	client      *http.Client

	Server *ServerService
	//DataCenter  *DataCenterService
}

func NewClient(username, password string) *Client {
	baseURL, err := url.Parse(defaultBaseURL)

	if err != nil {
		panic(err)
	}

	c := &Client{
		credentials: credentials{username, password},
		baseURL:     baseURL,
		client:      http.DefaultClient,
	}

	c.Server = &ServerService{client: c}
	//c.DataCenter = &DataCenterService{client: c}

	return c
}

//func NewClientWithToken(bearerToken, accountAlias string) Client {
//return Client{credentials: Credentials{bearerToken, accountAlias}}
//}

func (c *Client) Login() error {
	request, err := c.newRequest("POST", "authentication/login", c.credentials)

	if err != nil {
		return err
	}

	err = c.executeRequest(request, &c.user)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) newRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)

	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	buf := &bytes.Buffer{}
	if body != nil {
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")

	if c.user.BearerToken != "" {
		req.Header.Add("Authorization", "Bearer "+c.user.BearerToken)
	}

	return req, nil
}

func (c *Client) executeRequest(req *http.Request, v interface{}) error {
	resp, err := c.client.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; c < 200 || c > 299 {
		return fmt.Errorf("Server returns status %d", c)
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	return err
}
