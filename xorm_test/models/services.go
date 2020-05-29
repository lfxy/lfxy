package models

import (
	"time"
)

type Services struct {
	Id               int
	ServiceUuid      string
	ServiceBaseUuid  string
	UserId           string
	ProjectUuid      string
	CreateTime       time.Time
	EndTime          time.Time
	StorageType      string
	IfShared         int
	IfPersist        int
	IfJobPlan        int
	Status           string
	CpuAccumulate    float64
	MemoryAccumulate float64
	S3Accumulate     float64
	NfsAccumulate    float64
}
