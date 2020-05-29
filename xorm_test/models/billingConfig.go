package models

import ()

type BillingConfig struct {
	Id                int
	BillingConfigUuid string
	BillingPeriod     int
	StopServiceDelay  int
}
