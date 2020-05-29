package models

import ()

type ProjectMonth struct {
	Id                int
	ProjectMonthUuid  string
	Year              int
	Month             int
	ProjectUuid       string
	CpuUsage          int
	MemoryUsage       int
	StorageUsage      int
	StorageShareUsage int
	CpuBilling        float64
	MemoryBilling     float64
	StorageBilling    float64
	SumBilling        float64
}
