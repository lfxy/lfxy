package models

import (
	"time"
)

type ModuleMonth struct {
	Id              int
	ModuleMonthUuid string
	Year            int
	Month           int
	ModuleUuid      string
	CreateTime      time.Time
	EndTime         time.Time
	Cpu             int
	Memory          int
	Storage         int
	RunTime         int
	StorageTime     int
	StorageType     string
	StorageUseMode  string
	CpuBilling      float64
	MemoryBilling   float64
	StorageBilling  float64
	SumBilling      float64
}
