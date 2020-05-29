package models

import (
	"time"
)

type Instances struct {
	Id               int
	InstanceUuid     string
	InstanceName     string
	ServiceUuid      string
	CreateTime       time.Time
	EndTime          time.Time
	InstanceReplicas int
	StorageType      string
	Status           string
	CpuUsage         float32
	CpuLimit         float32
	MemoryUsage      float32
	MemoryLimit      float32
	S3Usage          float32
	S3Limit          float32
	NfsUsage         float32
	NfsLimit         float32
	UsedTime         time.Time
	CpuAccumulate    float64
	MemoryAccumulate float64
	S3Accumulate     float64
	NfsAccumulate    float64
}
