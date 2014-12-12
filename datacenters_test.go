package clcgo

import "testing"

func TestWorkingDataCentersURL(t *testing.T) {
	d := DataCenters{}
	d.SetCredentials(Credentials{AccountAlias: "AA"})
	u, err := d.URL()

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if e := ApiRoot + "/datacenters/AA"; u != e {
		t.Errorf("Expected URL to be '%s', was '%s'", e, u)
	}
}

// TODO: test unmarshalling of each entity !!!

func TestWorkingDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{DataCenter: DataCenter{ID: "abc123"}}
	d.SetCredentials(Credentials{AccountAlias: "AA"})
	u, err := d.URL()

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if e := ApiRoot + "/datacenters/AA/abc123/deploymentCapabilities"; u != e {
		t.Errorf("Expected URL to be '%s', was '%s'", e, u)
	}
}

func TestErroredDataCenterCapabilitiesURL(t *testing.T) {
	d := DataCenterCapabilities{}
	d.SetCredentials(Credentials{AccountAlias: "AA"})
	_, err := d.URL()

	if e := "Need a DataCenter with an ID"; err == nil || err.Error() != e {
		t.Errorf("Expected the error '%s', got '%v'", e, err)
	}
}

func TestSuccessfulDataCenterCapabilitiesUnmarshalling(t *testing.T) {
	d := DataCenterCapabilities{}
	j := []byte(`{"templates":[ { "name":"CENTOS-6-64-TEMPLATE" } ]}`)

	err := d.Unmarshal(j)

	if err != nil {
		t.Errorf("Expected no error, got '%s'", err)
	}

	if len(d.Templates) != 1 {
		t.Errorf("Expected Templates to have a length of 1")
	}

	if e := "CENTOS-6-64-TEMPLATE"; d.Templates[0].Name != e {
		t.Errorf("Expected '%s', but got '%s'", e, d.Templates[0].Name)
	}
}
