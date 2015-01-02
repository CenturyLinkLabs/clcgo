package clcgo

import (
	"encoding/json"
	"testing"

	"github.com/CenturyLinkLabs/clcgo/fakeapi"
	"github.com/stretchr/testify/assert"
)

func TestImplementations(t *testing.T) {
	es := []interface{}{
		new(DataCenterCapabilities),
		new(Server),
		new(Status),
	}
	for _, e := range es {
		assert.Implements(t, (*entity)(nil), e)
	}

	ses := []interface{}{
		new(Server),
		new(PublicIPAddress),
	}
	for _, se := range ses {
		assert.Implements(t, (*savableEntity)(nil), se)
	}
}

func TestServerJSONUnmarshalling(t *testing.T) {
	s := Server{}
	err := json.Unmarshal([]byte(fakeapi.ServerResponse), &s)

	assert.NoError(t, err)
	assert.Equal(t, "test-id", s.ID)
	assert.Equal(t, "Test Name", s.Name)
	assert.Equal(t, "123il", s.GroupID)
	assert.Len(t, s.Details.IPAddresses, 2)
	assert.Equal(t, "8.8.8.8", s.Details.IPAddresses[1].Public)
}

func TestWorkingServerURL(t *testing.T) {
	s := Server{ID: "abc123"}
	u, err := s.url("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/servers/AA/abc123", u)
}

func TestErroredServerURL(t *testing.T) {
	u, err := Server{}.url("AA")

	assert.EqualError(t, err, "An ID field is required to get a server")
	assert.Empty(t, u)
}

func TestURLMissingIDHavingUUID(t *testing.T) {
	u, err := Server{uuidURI: "/v2/alias/1234?uuid=true"}.url("AA")
	assert.NoError(t, err)
	assert.Equal(t, apiDomain+"/v2/alias/1234?uuid=true", u)
}

func TestServerRequestForSave(t *testing.T) {
	s := Server{
		Name:           "Test Name",
		GroupID:        "1234IL",
		SourceServerID: "TestID",
	}
	req, err := s.requestForSave("AA")
	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/servers/AA", req.URL)
	assert.Equal(t, s, req.Parameters)
}

func TestSuccessfulStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.statusFromResponse([]byte(fakeapi.ServerCreationSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/status", s.URI)
}

func TestErroredMissingStatusLinkStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.statusFromResponse([]byte(fakeapi.ServerCreationMissingStatusResponse))
	assert.Nil(t, s)
	assert.EqualError(t, err, "The creation response has no status link")
}

func TestSuccessfulIPAddressResponseForSave(t *testing.T) {
	s := Server{ID: "1234il"}
	ps := []Port{Port{Protocol: "TCP", Port: 31981}}
	i := PublicIPAddress{Server: s, Ports: ps}
	req, err := i.requestForSave("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiDomain+"/v2/servers/AA/1234il/publicIPAddresses", req.URL)
	assert.Equal(t, i, req.Parameters)
}

func TestErroredIPAddressResponseForSave(t *testing.T) {
	s := Server{}
	i := PublicIPAddress{Server: s}
	req, err := i.requestForSave("AA")

	assert.Equal(t, request{}, req)
	assert.EqualError(t, err, "A Server with an ID is required to add a Public IP Address")
}

func TestIPAddressStatusFromResponse(t *testing.T) {
	i := PublicIPAddress{}
	s, err := i.statusFromResponse([]byte(fakeapi.AddPublicIPAddressSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/status", s.URI)
}
