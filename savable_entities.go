package clcgo

import "encoding/json"

type request struct {
	URL        string
	Parameters interface{}
}

type savableEntity interface {
	requestForSave(string) (request, error)
}

type statusProvidingEntity interface {
	statusFromResponse([]byte) (*Status, error)
}

func SaveEntity(c Credentials, e savableEntity) (*Status, error) {
	return saveEntity(clcRequestor{}, c, e)
}

func saveEntity(r requestor, c Credentials, e savableEntity) (*Status, error) {
	req, err := e.requestForSave(c.AccountAlias)
	if err != nil {
		return nil, err
	}
	resp, err := r.PostJSON(c.BearerToken, req)
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
