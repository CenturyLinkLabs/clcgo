package clcgo

import (
	"errors"
	"fmt"
)

type DataCenters []DataCenter

type DataCenter struct {
	ID   string
	Name string
}

const DataCentersURL = APIRoot + "/datacenters/%s"

func (d DataCenters) URL(a string) (string, error) {
	return fmt.Sprintf(DataCentersURL, a), nil
}

type DataCenterCapabilities struct {
	DataCenter DataCenter `json:"-"`
	Templates  []struct {
		Name        string
		Description string
	}
}

const DataCenterCapabilitiesURL = DataCentersURL + "/%s/deploymentCapabilities"

func (d DataCenterCapabilities) URL(a string) (string, error) {
	if d.DataCenter.ID == "" {
		return "", errors.New("Need a DataCenter with an ID")
	}

	return fmt.Sprintf(DataCenterCapabilitiesURL, a, d.DataCenter.ID), nil
}
