package models

import (
	"time"
)

type ServicePlans struct {
	Id              int
	ServicePlanUuid string
	ServicePlanName string
	ServiceUuid     string
	CreateTime      time.Time
	DeleteTime      time.Time
	PeriodType      string
	StartTime       time.Time
	EndTime         time.Time
	Status          string
	IfActive        int
}
