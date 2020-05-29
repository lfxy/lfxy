package models

import ()

type WdfMonth struct {
	Id                int
	WdfMonthUuid      string
	Year              int
	Month             int
	CpuUsage          int
	MemoryUsage       int
	StorageUsage      int
	StorageShareUsage int
	S3Usage           int
	CpuBilling        float64
	MemoryBilling     float64
	StorageBilling    float64
	S3Billing         float64
	SumBillig         float64
}
