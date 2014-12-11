package clcgo

type Entity interface {
	URL() (string, error)
	Unmarshal([]byte) error
}

const ApiRoot = "https://api.tier3.com/v2"

//TODO implement GetEntity

func getEntity(r Requestor, t string, e Entity) error {
	//TODO: add bearer token support
	url, err := e.URL()
	if err != nil {
		return err
	}
	j, err := r.GetJSON(url)
	if err != nil {
		return err
	}
	return e.Unmarshal(j)
}
