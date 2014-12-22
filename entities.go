package clcgo

import "encoding/json"

const (
	apiDomain        = "https://api.tier3.com"
	apiRoot          = apiDomain + "/v2"
	successfulStatus = "succeeded"
)

// TODO Document this for this library's developers? Is there a flag to get it
// to not publish in GoDoc?
type entity interface {
	url(string) (string, error)
}

type Status struct {
	Status string
	URI    string
}

func (s Status) url(a string) (string, error) {
	return apiDomain + s.URI, nil
}

func (s Status) HasSucceeded() bool {
	return s.Status == successfulStatus
}

type Link struct {
	ID   string `json:"id"`
	Rel  string `json:"rel"`
	HRef string `json:"href"`
}

func GetEntity(c Credentials, e entity) error {
	return getEntity(&clcRequestor{}, c, e)
}

func getEntity(r requestor, c Credentials, e entity) error {
	url, err := e.url(c.AccountAlias)
	if err != nil {
		return err
	}
	j, err := r.GetJSON(c.BearerToken, request{URL: url})
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &e)
}
