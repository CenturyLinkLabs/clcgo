package clcgo

type Entity interface {
	URL() (string, error)
	Unmarshal([]byte) error
}

const ApiRoot = "https://api.tier3.com/v2"

func GetEntity(t string, e Entity) error {
	return getEntity(&CLCRequestor{}, t, e)
}

func getEntity(r Requestor, t string, e Entity) error {
	url, err := e.URL()
	if err != nil {
		return err
	}
	j, err := r.GetJSON(t, url)
	if err != nil {
		return err
	}
	return e.Unmarshal(j)
}
