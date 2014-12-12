package clcgo

import (
	"encoding/json"
	"errors"
	"fmt"
)

type DataCenters struct {
	Credentials Credentials `json:"-"`
	DataCenters []DataCenter
}

type DataCenter struct {
	ID   string
	Name string
}

const DataCentersURL = ApiRoot + "/datacenters/%s"

func (d *DataCenters) SetCredentials(c Credentials) {
	d.Credentials = c
}

func (d DataCenters) URL() (string, error) {
	return fmt.Sprintf(DataCentersURL, d.Credentials.AccountAlias), nil
}

func (d *DataCenters) Unmarshal(j []byte) error {
	return json.Unmarshal(j, &d.DataCenters)
}

type DataCenterCapabilities struct {
	DataCenter  DataCenter  `json:"-"`
	Credentials Credentials `json:"-"`
	Templates   []struct {
		Name        string
		Description string
	}
}

const DataCenterCapabilitiesURL = DataCentersURL + "/%s/deploymentCapabilities"

func (d *DataCenterCapabilities) SetCredentials(c Credentials) {
	d.Credentials = c
}

func (d DataCenterCapabilities) URL() (string, error) {
	if d.DataCenter.ID == "" {
		return "", errors.New("Need a DataCenter with an ID")
	}

	return fmt.Sprintf(DataCenterCapabilitiesURL, d.Credentials.AccountAlias, d.DataCenter.ID), nil
}

func (d *DataCenterCapabilities) Unmarshal(j []byte) error {
	return json.Unmarshal(j, &d)
}
