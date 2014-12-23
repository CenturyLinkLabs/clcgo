package clcgo

import "errors"

const authenticationURL = apiRoot + "/authentication/login"

type Credentials struct {
	BearerToken  string `json:"bearerToken"`
	AccountAlias string `json:"accountAlias"`
	Username     string `json:"username"` // TODO: nonexistant in get, extract to creation params?
	Password     string `json:"password"` // TODO: nonexistant in get, extract to creation params?
}

func (c Credentials) requestForSave(a string) (request, error) {
	return request{URL: authenticationURL, Parameters: c}, nil
}

func FetchCredentials(u string, p string) (Credentials, error) {
	return fetchCredentials(&clcRequestor{}, u, p)
}

func fetchCredentials(r requestor, u string, p string) (Credentials, error) {
	c := Credentials{Username: u, Password: p}
	_, err := saveEntity(r, c, &c)
	if err != nil {
		if rerr, ok := err.(RequestError); ok && rerr.StatusCode == 400 {
			err = errors.New("There was a problem with your credentials")
		}

		return Credentials{}, err
	}

	return c, nil
}
