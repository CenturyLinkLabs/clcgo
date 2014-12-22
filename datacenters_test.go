package clcgo

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorkingDataCentersURL(t *testing.T) {
	d := DataCenters{}
	u, err := d.url("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/datacenters/AA", u)
}

func TestDataCenterJSONUnmarshalling(t *testing.T) {
	j := `{"id": "foo", "name": "bar"}`

	dc := DataCenter{}
	err := json.Unmarshal([]byte(j), &dc)

	assert.NoError(t, err)

	assert.Equal(t, "foo", dc.ID)

	assert.Equal(t, "bar", dc.Name)
}

func TestWorkingDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{DataCenter: DataCenter{ID: "abc123"}}
	u, err := d.url("AA")

	assert.NoError(t, err)
	assert.Equal(t, apiRoot+"/datacenters/AA/abc123/deploymentCapabilities", u)
}

func TestErroredDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{}
	_, err := d.url("AA")
	assert.EqualError(t, err, "Need a DataCenter with an ID")
}

func TestSuccessfulDataCenterCapabilitiesUnmarshalling(t *testing.T) {
	d := DataCenterCapabilities{}
	err := json.Unmarshal([]byte(dataCenterCapabilitiesResponse), &d)

	assert.NoError(t, err)
	assert.Len(t, d.Templates, 1)
	assert.Equal(t, "Name", d.Templates[0].Name)
	assert.Equal(t, "Description", d.Templates[0].Description)
}
