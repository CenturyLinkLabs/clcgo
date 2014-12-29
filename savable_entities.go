package clcgo

type request struct {
	URL        string
	Parameters interface{}
}

// TODO Document this for this library's developers? Is a comment for an
// unexported interface suppressed from Godoc?
type savableEntity interface {
	requestForSave(string) (request, error)
}

// TODO Document this for this library's developers? Is a comment for an
// unexported interface suppressed from Godoc?
type statusProvidingEntity interface {
	statusFromResponse([]byte) (*Status, error)
}
