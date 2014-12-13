package clcgo

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestWorkingDataCentersURL(t *testing.T) {
	d := DataCenters{}
	u, err := d.URL("AA")

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if e := APIRoot + "/datacenters/AA"; u != e {
		t.Errorf("Expected URL to be '%s', was '%s'", e, u)
	}
}

func TestDataCenterJSONUnmarshalling(t *testing.T) {
	template := `{"id": "%s", "name": "%s"}`
	id := "foo"
	name := "bar"
	j := fmt.Sprintf(template, id, name)

	dc := DataCenter{}
	err := json.Unmarshal([]byte(j), &dc)

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if dc.ID != id {
		t.Errorf("Expected ID to be '%s', was '%s'", id, dc.ID)
	}

	if dc.Name != name {
		t.Errorf("Expected Name to be '%s', was '%s'", name, dc.Name)
	}
}

func TestWorkingDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{DataCenter: DataCenter{ID: "abc123"}}
	u, err := d.URL("AA")

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if e := APIRoot + "/datacenters/AA/abc123/deploymentCapabilities"; u != e {
		t.Errorf("Expected URL to be '%s', was '%s'", e, u)
	}
}

func TestErroredDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{}
	_, err := d.URL("AA")

	if e := "Need a DataCenter with an ID"; err == nil || err.Error() != e {
		t.Errorf("Expected the error '%s', got '%v'", e, err)
	}
}

func TestSuccessfulDataCenterCapabilitiesUnmarshalling(t *testing.T) {
	templates := `{"templates":[ { "name": "%s", "description": "%s" } ]}`
	name := "foo"
	description := "bar"
	j := fmt.Sprintf(templates, name, description)

	d := DataCenterCapabilities{}
	err := json.Unmarshal([]byte(j), &d)

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if len(d.Templates) != 1 {
		t.Errorf("Expected Templates to have a length of 1")
	}

	template := d.Templates[0]

	if template.Name != name {
		t.Errorf("Expected Name to be '%s', was '%s'", name, template.Name)
	}

	if template.Description != description {
		t.Errorf("Expected Description to be '%s', was '%s'", name, template.Description)
	}
}
