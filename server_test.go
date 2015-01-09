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
		new(Credentials),
		new(Status),
	}
	for _, e := range es {
		assert.Implements(t, (*Entity)(nil), e)
	}

	ses := []interface{}{
		new(Server),
		new(PublicIPAddress),
	}
	for _, se := range ses {
		assert.Implements(t, (*SavableEntity)(nil), se)
	}
}

func TestServerJSONUnmarshalling(t *testing.T) {
	s := Server{}
	err := json.Unmarshal([]byte(fakeapi.ServerResponse), &s)

	assert.NoError(t, err)
	assert.Equal(t, "test-id", s.ID)
	assert.Equal(t, "Test Name", s.Name)
	assert.Equal(t, "active", s.Status)
	assert.Equal(t, "123il", s.GroupID)
	assert.Len(t, s.Details.IPAddresses, 2)
	assert.Equal(t, "8.8.8.8", s.Details.IPAddresses[1].Public)
	assert.Equal(t, "started", s.Details.PowerState)
}

func TestIsActive(t *testing.T) {
	s := Server{}
	assert.False(t, s.IsActive())

	s.Status = "active"
	s.Details.PowerState = "paused"
	assert.False(t, s.IsActive())

	s.Details.PowerState = "started"
	assert.True(t, s.IsActive())
}

func TestIsPaused(t *testing.T) {
	s := Server{}
	assert.False(t, s.IsPaused())

	s.Details.PowerState = "started"
	assert.False(t, s.IsPaused())

	s.Details.PowerState = "paused"
	assert.True(t, s.IsPaused())
}

func TestSuccessfulServerURL(t *testing.T) {
	s := Server{ID: "abc123"}
	u, err := s.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/servers/AA/abc123", u)
}

func TestErroredServerURL(t *testing.T) {
	u, err := Server{}.URL("AA")

	assert.EqualError(t, err, "An ID field is required to get a server")
	assert.Empty(t, u)
}

func TestURLMissingIDHavingUUID(t *testing.T) {
	u, err := Server{uuidURI: "/v2/alias/1234?uuid=true"}.URL("AA")
	assert.NoError(t, err)
	assert.Equal(t, apiDomain+"/v2/alias/1234?uuid=true", u)
}

func TestServerRequestForSave(t *testing.T) {
	s := Server{
		Name:           "Test Name",
		GroupID:        "1234IL",
		SourceServerID: "TestID",
	}
	req, err := s.RequestForSave("AA")
	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/servers/AA", req.URL)
	assert.Equal(t, s, req.Parameters)
}

func TestSuccessfulStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.StatusFromResponse([]byte(fakeapi.ServerCreationSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/v2/operations/alias/status/test-status-id", s.URI)
}

func TestErroredMissingStatusLinkStatusFromResponse(t *testing.T) {
	srv := Server{}
	s, err := srv.StatusFromResponse([]byte(fakeapi.ServerCreationMissingStatusResponse))
	assert.Nil(t, s)
	assert.EqualError(t, err, "The creation response has no status link")
}

func TestSuccessfulIPAddressResponseForSave(t *testing.T) {
	s := Server{ID: "1234il"}
	ps := []Port{Port{Protocol: "TCP", Port: 31981}}
	i := PublicIPAddress{Server: s, Ports: ps}
	req, err := i.RequestForSave("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiDomain+"/v2/servers/AA/1234il/publicIPAddresses", req.URL)
	assert.Equal(t, i, req.Parameters)
}

func TestErroredIPAddressResponseForSave(t *testing.T) {
	s := Server{}
	i := PublicIPAddress{Server: s}
	req, err := i.RequestForSave("AA")

	assert.Equal(t, request{}, req)
	assert.EqualError(t, err, "A Server with an ID is required to add a Public IP Address")
}

func TestIPAddressStatusFromResponse(t *testing.T) {
	i := PublicIPAddress{}
	s, err := i.StatusFromResponse([]byte(fakeapi.AddPublicIPAddressSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/status", s.URI)
}

func TestCredentialsJSONUnmarshalling(t *testing.T) {
	c := Credentials{}
	err := json.Unmarshal([]byte(fakeapi.ServerCredentialsResponse), &c)

	assert.NoError(t, err)
	assert.Equal(t, "root", c.Username)
	assert.Equal(t, "p4ssw0rd", c.Password)
}

func TestSuccessfulCredentialsURL(t *testing.T) {
	s := Server{ID: "abc123"}
	c := Credentials{Server: s}
	u, err := c.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/servers/AA/abc123/credentials", u)
}

func TestErroredCredentialsURL(t *testing.T) {
	c := Credentials{}
	c.Server = Server{ID: ""}

	u, err := c.URL("AA")
	assert.EqualError(t, err, "A Server with an ID is required to fetch credentials")
	assert.Empty(t, u)
}
