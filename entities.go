package clcgo

import "encoding/json"

type Entity interface {
	URL(string) (string, error)
}

type Status struct {
	Status string
	URI    string
}

func (s Status) URL(a string) (string, error) {
	return APIDomain + s.URI, nil
}

func (s Status) HasSucceeded() bool {
	return s.Status == "succeeded"
}

type Link struct {
	ID   string `json:"id"`
	Rel  string `json:"rel"`
	HRef string `json:"href"`
}

const (
	APIDomain = "https://api.tier3.com"
	APIRoot   = APIDomain + "/v2"
)

func GetEntity(c Credentials, e Entity) error {
	return getEntity(&CLCRequestor{}, c, e)
}

func getEntity(r Requestor, c Credentials, e Entity) error {
	url, err := e.URL(c.AccountAlias)
	if err != nil {
		return err
	}
	j, err := r.GetJSON(c.BearerToken, url)
	if err != nil {
		return err
	}

	return json.Unmarshal(j, &e)
}
