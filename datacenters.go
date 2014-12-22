package clcgo

import (
	"errors"
	"fmt"
)

const (
	dataCentersURL            = apiRoot + "/datacenters/%s"
	dataCenterCapabilitiesURL = dataCentersURL + "/%s/deploymentCapabilities"
)

type DataCenters []DataCenter

type DataCenter struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type DataCenterCapabilities struct {
	DataCenter DataCenter `json:"-"`
	Templates  []struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
}

func (d DataCenters) url(a string) (string, error) {
	return fmt.Sprintf(dataCentersURL, a), nil
}

func (d DataCenterCapabilities) url(a string) (string, error) {
	if d.DataCenter.ID == "" {
		return "", errors.New("Need a DataCenter with an ID")
	}

	return fmt.Sprintf(dataCenterCapabilitiesURL, a, d.DataCenter.ID), nil
}
