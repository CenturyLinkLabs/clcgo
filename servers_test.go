package clcgo

import (
	"testing"

	"github.com/CenturyLinkLabs/clcgo/fakeapi"
	"github.com/stretchr/testify/assert"
)

func TestSuccessfulPauseServersRequestForSave(t *testing.T) {
	s := Server{ID: "test-id"}
	p := PauseServer{Server: s}
	req, err := p.RequestForSave("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/operations/AA/servers/pause", req.URL)

	sids, ok := req.Parameters.([]string)
	if assert.True(t, ok) {
		assert.Len(t, sids, 1)
		assert.Equal(t, "test-id", sids[0])
	}
}

func TestErroredPauseServersRequestForSave(t *testing.T) {
	p := PauseServer{}
	req, err := p.RequestForSave("AA")

	assert.Equal(t, request{}, req)
	assert.EqualError(t, err, "PauseServer requires a Server with an ID to pause!")
}

func TestPauseServersStatusFromResponse(t *testing.T) {
	s := Server{ID: "test-id"}
	p := PauseServer{Server: s}
	st, err := p.StatusFromResponse([]byte(fakeapi.PauseServersSuccessfulResponse))
	assert.NoError(t, err)
	assert.Equal(t, "/path/to/status", st.URI)
}
