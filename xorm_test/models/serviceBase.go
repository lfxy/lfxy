package models

import ()

type ServiceBase struct {
	Id                    int
	ServiceBaseUuid       string
	ServiceName           string
	Version               string
	AssociatedServiceUuid string
}
