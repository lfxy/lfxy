package models

import ()

type ServiceMonth struct {
	Id                int
	ServiceMonthUuid  string
	Year              int
	Month             int
	ServiceUuid       string
	CpuUsage          int
	MemoryUsage       int
	StorageUsage      int
	StorageShareUsage int
	CpuBilling        float64
	MemoryBilling     float64
	StorageBilling    float64
	SumBilling        float64
}
