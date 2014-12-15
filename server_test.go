package clcgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkingServerURL(t *testing.T) {
	s := Server{ID: "abc123"}
	u, err := s.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, APIRoot+"/servers/AA/abc123", u)
}

func TestErroredServerURL(t *testing.T) {
	s := Server{}
	u, err := s.URL("AA")

	assert.EqualError(t, err, "The server needs an ID attribute to generate a URL")
	assert.Empty(t, u)
}
