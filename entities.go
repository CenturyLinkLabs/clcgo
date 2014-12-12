package clcgo

type Entity interface {
	URL(string) (string, error)
	Unmarshal([]byte) error
}

const APIRoot = "https://api.tier3.com/v2"

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
	return e.Unmarshal(j)
}
