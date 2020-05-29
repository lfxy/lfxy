package models

import ()

type UserMonth struct {
	Id                int
	UserMonthUuid     string
	Year              int
	Month             int
	UserId            string
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
