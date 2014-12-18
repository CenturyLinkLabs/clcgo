package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requestor interface {
	PostJSON(string, Request) ([]byte, error)
	GetJSON(string, Request) ([]byte, error)
}

type CLCRequestor struct{}

type RequestError struct {
	Err        string
	StatusCode int
}

func (r RequestError) Error() string {
	return r.Err
}

func (r CLCRequestor) PostJSON(t string, req Request) ([]byte, error) {
	json, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	hr, err := http.NewRequest("POST", req.URL, strings.NewReader(string(json)))
	if err != nil {
		return nil, err
	}

	if t != "" {
		hr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
	}
	hr.Header.Add("Content-Type", "application/json")
	hr.Header.Add("Accepts", "application/json")

	resp, err := client.Do(hr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200, 201, 202:
		return body, nil
	case 401:
		return body, RequestError{"Your bearer token was rejected", 401}
	default:
		return body, RequestError{"Got an unexpected status code", resp.StatusCode}
	}
}

func (r CLCRequestor) GetJSON(t string, req Request) ([]byte, error) {
	client := http.Client{}

	hr, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		return nil, err
	}

	hr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
	hr.Header.Add("Content-Type", "application/json")
	hr.Header.Add("Accepts", "application/json")

	resp, err := client.Do(hr)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	switch resp.StatusCode {
	case 200:
		return body, nil
	case 401:
		return body, RequestError{"Your bearer token was rejected", 401}
	default:
		return body, RequestError{"Got an unexpected status code", resp.StatusCode}
	}
}

func typeFromLinks(t string, ls []Link) (Link, error) {
	for _, l := range ls {
		if l.Rel == t {
			return l, nil
		}
	}

	return Link{}, fmt.Errorf("No link of type '%s'", t)
}
