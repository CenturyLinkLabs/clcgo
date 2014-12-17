package clcgo

type SavableEntity interface {
	URLForSave(string) (string, error)
	ParametersForSave() (interface{}, error)
	StatusFromResponse([]byte) (*Status, error)
}

func SaveEntity(c Credentials, e SavableEntity) (*Status, error) {
	return saveEntity(CLCRequestor{}, c, e)
}

func saveEntity(r Requestor, c Credentials, e SavableEntity) (*Status, error) {
	url, err := e.URLForSave(c.AccountAlias)
	if err != nil {
		return nil, err
	}
	params, err := e.ParametersForSave()
	if err != nil {
		return nil, err
	}
	resp, err := r.PostJSON(c.BearerToken, url, params)
	if err != nil {
		return nil, err
	}
	status, err := e.StatusFromResponse(resp)
	if err != nil {
		return nil, err
	}

	return status, nil
}
