package clcgo

import (
	"errors"
	"fmt"
)

const (
	groupURL = apiRoot + "/groups/%s/%s"
)

// A Group resource can be used to discover the hierarchy of the created groups
// for your account for a given datacenter. The Group's ID can either be for
// one you have created, or for a datacenter's hardware group (which can be
// determined via the DataCenterGroup resource).
type Group struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Type   string  `json:"type"`
	Groups []Group `json:"groups"`
}

func (g Group) URL(a string) (string, error) {
	if g.ID == "" {
		return "", errors.New("An ID field is required to get a group")
	}

	return fmt.Sprintf(groupURL, a, g.ID), nil
}
