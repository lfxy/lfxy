package models

import ()

type StoragePrice struct {
	Id               int
	StoragePriceUuid string
	StorageName      string
	BillingMode      string
	FreeValue        int
	UnitValue        float64
}
