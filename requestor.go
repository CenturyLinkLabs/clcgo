package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var DefaultHTTPClient = &http.Client{}

type requestor interface {
	PostJSON(string, request) ([]byte, error)
	GetJSON(string, request) ([]byte, error)
}

type clcRequestor struct{}

type modelStates map[string][]string

type RequestError struct {
	Message    string
	StatusCode int
	Errors     modelStates
}

type invalidReqestResponse struct {
	Message    string      `json:"message"`
	ModelState modelStates `json:"modelState"`
}

func (r RequestError) Error() string {
	return r.Message
}

func (r clcRequestor) PostJSON(t string, req request) ([]byte, error) {
	j, err := json.Marshal(req.Parameters)
	if err != nil {
		return nil, err
	}

	hr, err := http.NewRequest("POST", req.URL, strings.NewReader(string(j)))
	if err != nil {
		return nil, err
	}

	if t != "" {
		hr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
	}
	hr.Header.Add("Content-Type", "application/json")
	hr.Header.Add("Accepts", "application/json")

	resp, err := DefaultHTTPClient.Do(hr)
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
	case 400:
		var e invalidReqestResponse
		err := json.Unmarshal(body, &e)
		if err != nil {
			return body, err
		}

		return body, RequestError{Message: e.Message, StatusCode: 400, Errors: e.ModelState}
	case 401:
		return body, RequestError{Message: "Your bearer token was rejected", StatusCode: 401}
	default:
		return body, RequestError{Message: "Got an unexpected status code", StatusCode: resp.StatusCode}
	}
}

func (r clcRequestor) GetJSON(t string, req request) ([]byte, error) {
	hr, err := http.NewRequest("GET", req.URL, nil)
	if err != nil {
		return nil, err
	}

	hr.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
	hr.Header.Add("Content-Type", "application/json")
	hr.Header.Add("Accepts", "application/json")

	resp, err := DefaultHTTPClient.Do(hr)
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
		return body, RequestError{Message: "Your bearer token was rejected", StatusCode: 401}
	default:
		return body, RequestError{Message: "Got an unexpected status code", StatusCode: resp.StatusCode}
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
