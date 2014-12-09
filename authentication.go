package clcgo

import "encoding/json"

type Credentials struct {
	BearerToken string
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
		return Credentials{}, err
	} else {
		var credentials Credentials
		json.Unmarshal(response, &credentials)
		return credentials, nil
	}
}
