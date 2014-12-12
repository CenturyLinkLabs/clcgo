package clcgo

type Entity interface {
	SetCredentials(Credentials)
	URL() (string, error)
	Unmarshal([]byte) error
}

const ApiRoot = "https://api.tier3.com/v2"

func GetEntity(c Credentials, e Entity) error {
	return getEntity(&CLCRequestor{}, c, e)
}

func getEntity(r Requestor, c Credentials, e Entity) error {
	e.SetCredentials(c)
	url, err := e.URL()
	if err != nil {
		return err
	}
	j, err := r.GetJSON(c.BearerToken, url)
	if err != nil {
		return err
	}
	return e.Unmarshal(j)
}
