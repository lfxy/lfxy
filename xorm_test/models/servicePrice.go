package models

import ()

type ServicePrice struct {
	Id               int
	ServicePriceUuid string
	ResourceName     string
	BillingMode      string
	FreeValue        int
	UnitValue        float64
}
