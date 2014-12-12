package clcgo

import (
	"encoding/json"
	"errors"
)

type Credentials struct {
	BearerToken  string
	AccountAlias string
}

type authParameters struct {
	Username string
	Password string
}

const AuthenticationURL = "https://api.tier3.com/v2/authentication/login"

func FetchCredentials(username string, password string) (Credentials, error) {
	return fetchCredentials(&CLCRequestor{}, username, password)
}

func fetchCredentials(client Requestor, username string, password string) (Credentials, error) {
	c := authParameters{username, password}
	response, err := client.PostJSON(AuthenticationURL, c)

	if err != nil {
		if rerr, ok := err.(RequestError); ok && rerr.StatusCode == 400 {
			err = errors.New("There was a problem with your credentials")
		}

		return Credentials{}, err
	}

	var credentials Credentials
	json.Unmarshal(response, &credentials)
	return credentials, nil
}
