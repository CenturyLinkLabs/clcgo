package clcgo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

type Requestor interface {
	PostJSON(url string, v interface{}) ([]byte, error)
}

type CLCRequestor struct{}

func (r *CLCRequestor) PostJSON(url string, v interface{}) ([]byte, error) {
	json, _ := json.Marshal(v)

	client := &http.Client{}
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

	// TODO: Check for unexpected response codes

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
