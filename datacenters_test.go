package clcgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkingDataCentersURL(t *testing.T) {
	d := DataCenters{}
	u, err := d.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, APIRoot+"/datacenters/AA", u)
}

// TODO: test unmarshalling of each entity !!!

func TestWorkingDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{DataCenter: DataCenter{ID: "abc123"}}
	u, err := d.URL("AA")

	assert.NoError(t, err)
	assert.Equal(t, APIRoot+"/datacenters/AA/abc123/deploymentCapabilities", u)
}

func TestErroredDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{}
	_, err := d.URL("AA")
	assert.EqualError(t, err, "Need a DataCenter with an ID")
}

func TestSuccessfulDataCenterCapabilitiesUnmarshalling(t *testing.T) {
	d := DataCenterCapabilities{}
	j := []byte(`{"templates":[ { "name":"CENTOS-6-64-TEMPLATE" } ]}`)

	err := d.Unmarshal(j)

	assert.NoError(t, err)
	assert.Len(t, d.Templates, 1)
	assert.Equal(t, "CENTOS-6-64-TEMPLATE", d.Templates[0].Name)
}
