package clcgo

const (
	apiDomain        = "https://api.tier3.com"
	apiRoot          = apiDomain + "/v2"
	successfulStatus = "succeeded"
)

// TODO Document this for this library's developers? Is a comment for an
// unexported interface suppressed from Godoc?
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
