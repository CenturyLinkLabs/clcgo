package clcgo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requestor interface {
	PostJSON(url string, v interface{}) ([]byte, error)
	GetJSON(t string, url string) ([]byte, error)
}

type CLCRequestor struct{}

type RequestError struct {
	Err        string
	StatusCode int
}

func (r RequestError) Error() string {
	return r.Err
}

func (r CLCRequestor) PostJSON(url string, v interface{}) ([]byte, error) {
	json, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}

	client := http.Client{}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(json)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, RequestError{"Got an unexpected status code", resp.StatusCode}
	}

	return body, nil
}

func (r CLCRequestor) GetJSON(t string, url string) ([]byte, error) {
	// TODO: DisableKeepAlives is pretty much just for tests...
	client := http.Client{Transport: &http.Transport{DisableKeepAlives: true}}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", t))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return body, RequestError{"Got an unexpected status code", resp.StatusCode}
	}

	return body, nil
}
