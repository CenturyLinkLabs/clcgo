package clcgo

type Request struct {
	URL        string
	Parameters interface{}
}

type SavableEntity interface {
	RequestForSave(string) (Request, error)
	StatusFromResponse([]byte) (*Status, error)
}

func SaveEntity(c Credentials, e SavableEntity) (*Status, error) {
	return saveEntity(CLCRequestor{}, c, e)
}

func saveEntity(r Requestor, c Credentials, e SavableEntity) (*Status, error) {
	req, err := e.RequestForSave(c.AccountAlias)
	if err != nil {
		return nil, err
	}
	resp, err := r.PostJSON(c.BearerToken, req)
	if err != nil {
		return nil, err
	}
	status, err := e.StatusFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return status, nil
}
