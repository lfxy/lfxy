package models

import (
	"time"
)

type StoragePurchase struct {
	Id                  int
	StoragePurchaseUuid string
	UserId              string
	StoragePriceUuid    string
	Value               int
	Days                string
	CreateTime          time.Time
	EndTime             time.Time
	UsedDays            int
	CostBilling         float64
}
