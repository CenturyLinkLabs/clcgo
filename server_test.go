package clcgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServerJSONUnmarshalling(t *testing.T) {
	j := `{"id": "foo", "name": "bar", "groupId": "123il"}`
	s := Server{}
	err := json.Unmarshal([]byte(j), &s)

	assert.NoError(t, err)
	assert.Equal(t, "foo", s.ID)
	assert.Equal(t, "bar", s.Name)
	assert.Equal(t, "123il", s.GroupID)
}

func TestWorkingServerURL(t *testing.T) {
	s := Server{ID: "abc123"}
	u, err := s.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, APIRoot+"/servers/AA/abc123", u)
}

func TestErroredServerURL(t *testing.T) {
	u, err := Server{}.URL("AA")

	assert.EqualError(t, err, "An ID field is required to get a server")
	assert.Empty(t, u)
}

func TestURLForSave(t *testing.T) {
	url, err := Server{}.URLForSave("AA")
	assert.NoError(t, err)
	assert.Equal(t, APIRoot+"/servers/AA", url)
}

func TestSuccessfulParametersForSave(t *testing.T) {
	s := Server{
		Name:           "Test Name",
		GroupID:        "1234IL",
		SourceServerID: "TestID",
	}
	p, err := s.ParametersForSave()
	assert.NoError(t, err)
	assert.Equal(t, s, p)
}

func TestUnsuccessfulParametersForSave(t *testing.T) {
	p, err := Server{}.ParametersForSave()
	assert.Nil(t, p)
	assert.EqualError(t, err, "The following fields are required to save: Name, GroupID, SourceServerID")
}

func TestSuccessfulStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.StatusFromResponse([]byte(serverCreationSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/status", s.URI)
}

func TestErroredMissingStatusLinkStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.StatusFromResponse([]byte(serverCreationMissingStatusResponse))
	assert.Nil(t, s)
	assert.EqualError(t, err, "The creation response has no status link")
}
