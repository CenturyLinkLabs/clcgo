package clcgo

import (
	"errors"
	"fmt"
)

const (
	dataCentersURL            = apiRoot + "/datacenters/%s"
	dataCenterCapabilitiesURL = dataCentersURL + "/%s/deploymentCapabilities"
)

// The DataCenters resource can retrieve a list of available DataCenters.
type DataCenters []DataCenter

// A DataCenter resource can either be returned by the DataCenters resource, or
// built manually. It should be used in conjunction with the
// DataCenterCapabilities resource to request information about it.
//
// You must supply the ID if you are building this object manually.
type DataCenter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// DataCenterCapabilities gets more information about a specific DataCenter.
// You must supply the associated DataCenter object.
type DataCenterCapabilities struct {
	DataCenter         DataCenter `json:"-"`
	DeployableNetworks []struct {
		Name      string `json:"name"`
		NetworkID string `json:"networkId"`
		Type      string `json:"type"`
		AccountID string `json:"accountID"`
	} `json:"deployableNetworks"`
	Templates []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"templates"`
}

func (d DataCenters) URL(a string) (string, error) {
	return fmt.Sprintf(dataCentersURL, a), nil
}

func (d DataCenterCapabilities) URL(a string) (string, error) {
	if d.DataCenter.ID == "" {
		return "", errors.New("Need a DataCenter with an ID")
	}

	return fmt.Sprintf(dataCenterCapabilitiesURL, a, d.DataCenter.ID), nil
}
