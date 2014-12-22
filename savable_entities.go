package clcgo

type request struct {
	URL        string
	Parameters interface{}
}

type savableEntity interface {
	requestForSave(string) (request, error)
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
	status, err := e.statusFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return status, nil
}
