package models

import (
	"time"
)

type ModuleHistory struct {
	Id                int
	ModuleHistoryUuid string
	ModuleUuid        string
	ModuleName        string
	ServiceUuid       string
	ServiceName       string
	UserId            string
	ProjectUuid       string
	ProjectName       string
	CreateTime        time.Time
	EndTime           time.Time
	Cpu               int
	Memory            int
	Storage           int
	RunTime           int
	StorageTime       int
	StorageType       string
	StorageUseMode    string
	CpuBilling        float64
	MemoryBilling     float64
	StorageBilling    float64
	SumBilling        int
}
