package models

import (
	"time"
)

type Projects struct {
	Id               int
	ProjectUuid      string
	ProjectName      string
	UserId           string
	ProjectType      string
	BillingMode      string
	CreateTime       time.Time
	EndTime          time.Time
	Status           string
	CpuAccumulate    float64
	MemoryAccumulate float64
	S3Accumulate     float64
	NfsAccumulate    float64
}
